pipeline {
    agent any
    
    environment {
        IMAGE_NAME = 'argyarijal/mqtt-client-ta:latest'
        CONTAINER_NAME = 'service-mqtt-golang'
        BRANCH_NAME = "main"
        MSG_COMMIT = sh(script: "git log -1 --pretty=%B ${env.GIT_COMMIT}", returnStdout: true).trim()
    }
    
    stages {
       
        stage('Skip pattern') {
            steps {
                script {
                    if (MSG_COMMIT == 'generated by jenkins') {
                        echo 'Build Canceled'
                        currentBuild.result = 'ABORTED'
                        throw new Exception()
                    }
                }
            }
        }

        stage('Docker Build and Push') {
            steps {
                withDockerRegistry([credentialsId: "docker-creds", url: ""]) {
                    retry(3) {
                        timeout(time: 25, unit: 'MINUTES') {
                            sh 'printenv'
                            sh 'DOCKER_BUILDKIT=1 docker build --rm=false -t ${IMAGE_NAME} .'
                            sh 'docker push ${IMAGE_NAME}'
                        }
                    }
                }
            }
        }

        stage('Prune Docker Data') {
            steps {
                sh '''
                if [ $(docker ps -a -q -f name=${CONTAINER_NAME}) ]; then
                    docker stop ${CONTAINER_NAME} || true
                    docker rm ${CONTAINER_NAME} || true
                fi
                '''
            }
        }
        stage('Deploy to portainer') {
            steps {
                script {
                sh "curl -k -X POST https://193.203.167.97:9443/api/stacks/webhooks/162ce309-938b-441c-a044-450dad4460e3"
                echo "Deployed to Portainer"
                }
            }
        }
    }

    post {
        success {
            echo 'Build and deployment successful!'
            discordSend (
                webhookURL: "https://discord.com/api/webhooks/1244675279239254077/OBqglM-TvJsGJ-PJdxxthw0-mkEzkFPJb4phGpDIBvd0jSpnJ_HyZeNc6C8ML3lJnR9Y",
                title: "${JOB_NAME}",
                description: "Build Success \n - Link Build : ${env.BUILD_URL} \n - Image : `${env.IMAGE_NAME}` \n - Branch Name : `${env.BRANCH_NAME}` \n - Commit MSG : `${MSG_COMMIT}`",
                result: currentBuild.currentResult.toString(),
                footer: "Footer Text"
            )
        }

        failure {
            echo 'Build or deployment failed!'
            discordSend (
                webhookURL: "https://discord.com/api/webhooks/1244675279239254077/OBqglM-TvJsGJ-PJdxxthw0-mkEzkFPJb4phGpDIBvd0jSpnJ_HyZeNc6C8ML3lJnR9Y",
                title: "${JOB_NAME}",
                description: "Build Failed \n - Link Build : ${env.BUILD_URL} \n - Image : `${env.IMAGE_NAME}` \n - Branch Name : `${env.BRANCH_NAME}` \n - Commit MSG : `${MSG_COMMIT}`",
                result: currentBuild.currentResult.toString(),
                footer: "Footer Text"
            )
        }
    }
}
