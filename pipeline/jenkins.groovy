pipeline {
    agent any

    tools {
        go 'go1.26.5'
    }

    parameters {
        choice(
            name: 'OS',
            choices: ['linux', 'darwin', 'windows'],
            description: 'Target operating system'
        )

        choice(
            name: 'ARCH',
            choices: ['amd64', 'arm64'],
            description: 'Target architecture'
        )

        booleanParam(
            name: 'SKIP_TESTS',
            defaultValue: false,
            description: 'Skip running tests'
        )

        booleanParam(
            name: 'SKIP_LINT',
            defaultValue: false,
            description: 'Skip source formatting check'
        )
    }

    options {
        timestamps()
        disableConcurrentBuilds()
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Environment') {
            steps {
                sh 'go version'
                sh 'git --version'
                echo "Target platform: ${params.OS}/${params.ARCH}"
            }
        }

        stage('Lint') {
            when {
                expression { !params.SKIP_LINT }
            }
            steps {
                sh '''
                    UNFORMATTED="$(find . -type f -name '*.go' -not -path './.git/*' -exec gofmt -l {} +)"

                    if [ -n "$UNFORMATTED" ]; then
                        echo "Files require gofmt:"
                        echo "$UNFORMATTED"
                        exit 1
                    fi

                    echo "Go formatting check passed"
                '''
            }
        }

        stage('Test') {
            when {
                expression { !params.SKIP_TESTS }
            }
            steps {
                sh 'go test -v ./...'
            }
        }

        stage('Build') {
            steps {
                script {
                    def extension = params.OS == 'windows' ? '.exe' : ''
                    def artifact = "bin/kbot-${params.OS}-${params.ARCH}${extension}"

                    sh '''
                        mkdir -p bin
                    '''

                    sh """
                        CGO_ENABLED=0 \\
                        GOOS=${params.OS} \\
                        GOARCH=${params.ARCH} \\
                        go build -trimpath -ldflags="-s -w" -o ${artifact} .
                    """

                    echo "Artifact: ${artifact}"
                }
            }
        }
    }

    post {
        success {
            echo "Build completed successfully for ${params.OS}/${params.ARCH}"
        }

        failure {
            echo "Build failed for ${params.OS}/${params.ARCH}"
        }

        always {
            archiveArtifacts artifacts: 'bin/*', allowEmptyArchive: true
        }
    }
}
