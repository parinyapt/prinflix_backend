/* 
Pipeline: Prinflix Backend Golang
Created by: Parinya Termkasipanich
Date: 24-11-2023

Required Jenkins Plugins:
- SSH Agent
- Generic Webhook Trigger Plugin
*/


pipeline {
    agent any

    stages {
        stage('Checkout') {
            steps {
                echo '--- Start Checkout Stage ---'
                checkout scm
                echo '--- Finish Checkout Stage ---'
            }
        }
        stage('Info') {
            steps {
                echo '--- Start Info Stage ---'
                echo "Git Commit: ${GIT_COMMIT}"
                echo "Git Tag: ${PAYLOAD_TAGNAME}"
                echo "Trigger Event : ${x_github_event}"
                echo "Trigger Action : ${PAYLOAD_ACTION}"
                echo '--- Finish Info Stage ---'
            }
        }
        stage('Check Trigger Event') {
            when {
                not {
                    anyOf {
                        allOf {
                            environment name: 'x_github_event', value: 'release'
                            environment name: 'PAYLOAD_ACTION', value: 'released'
                        }
                        allOf {
                            environment name: 'x_github_event', value: 'release'
                            environment name: 'PAYLOAD_ACTION', value: 'created'
                        }
                        allOf {
                            environment name: 'x_github_event', value: 'release'
                            environment name: 'PAYLOAD_ACTION', value: 'published'
                        }
                    }
                }
            }
            steps {
                script {
                    currentBuild.result = 'ABORTED'
                    error("Aborting Build Process")
                }
            }
        }
        stage('Build') {
            environment {
                IMAGE_REGISTRY = "ghcr.io/parinyapt"
                IMAGE_NAME = "prinflix-backend"
            }
            steps {
                echo '--- Start Build Stage ---'
                withCredentials([usernamePassword(
                    credentialsId: 'parinyapt-github-credentials',
                    usernameVariable: 'registry_username',
                    passwordVariable: 'registry_password'
                )]) {
                    sh 'docker login ghcr.io -u $registry_username -p $registry_password'
                }
                sh "docker build -t ${IMAGE_NAME}:${PAYLOAD_TAGNAME} --platform linux/amd64 ."
                sh "docker tag ${IMAGE_NAME}:${PAYLOAD_TAGNAME} ${IMAGE_REGISTRY}/${IMAGE_NAME}:${PAYLOAD_TAGNAME}"
                sh "docker push ${IMAGE_REGISTRY}/${IMAGE_NAME}:${PAYLOAD_TAGNAME}"
                sh "docker image rm ${IMAGE_NAME}:${PAYLOAD_TAGNAME}"
                sh "docker image rm ${IMAGE_REGISTRY}/${IMAGE_NAME}:${PAYLOAD_TAGNAME}"
            }
        }
        stage('Deploy') {
            environment {
                SSH_SCRIPT = """
                    cd ${SERVER_SERVICE_INFRASTRUCTURE_CODE_PATH}/prinflix/backend && \
                    IMAGE_TAG=${PAYLOAD_TAGNAME} docker-compose up -d
                """
            }
            steps {
                echo '--- Start Deploy Stage ---'
                sshagent(credentials:['ssh-prin-server-1']){
                    sh "ssh  -o StrictHostKeyChecking=no  ${SERVER_USERNAME}@${SERVER_IPADDRESS} -p ${SERVER_PORT} '${SSH_SCRIPT}'"
                }
                echo '--- Finish Deploy Stage ---'
            }
        }
    }
    post { 
        always { 
            sh """
            curl -X POST http://jenkins-notify-api/prinflix-backend-golang/v1/notify \
            -H 'AccessKey: 609e932f-ba4c-4657-8113-e483956e7052' \
            -H 'Content-Type: application/json' \
            -d '{"message":"Prinflix Backend : Deploy ${currentBuild.result}\\nDetail\\n- Git Commit: ${GIT_COMMIT}\\n- Git Tag: ${PAYLOAD_TAGNAME}\\n- Trigger Event : ${x_github_event}\\n- Trigger Action : ${PAYLOAD_ACTION}\\n- Build ID: ${BUILD_ID}\\n- Build Number: ${BUILD_NUMBER}\\n- Build Tag: ${BUILD_TAG}\\n- Build URL: ${BUILD_URL}\\n"}'
            """
        }
    }
}