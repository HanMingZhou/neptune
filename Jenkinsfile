pipeline {
    agent any

    environment {
        GIT_URL = 'https://github.com/HanMingZhou/neptune.git'
        GIT_BRANCH = 'master'

        HARBOR_REGISTRY = 'dockerhub.kubekey.local'
        HARBOR_PROJECT = 'neptune'
        HARBOR_CREDENTIAL_ID = 'harbor'

        DOCKER_BUILDKIT = '0'
    }

    options {
        timestamps()
        disableConcurrentBuilds()
    }

    stages {
        stage('网络诊断 (Network Check)') {
            steps {
                sh '''
                    set -e
                    echo "========== DNS 解析 =========="
                    nslookup github.com || echo "❌ DNS 解析 github.com 失败"

                    echo ""
                    echo "========== TCP 连通性 =========="
                    timeout 5 bash -c 'cat < /dev/null > /dev/tcp/github.com/443' \
                        && echo "✅ github.com:443 可达" \
                        || echo "❌ github.com:443 不可达"

                    echo ""
                    echo "========== HTTPS 连通性 =========="
                    curl -sS --connect-timeout 5 --max-time 10 -o /dev/null -w "HTTP Status: %{http_code}\\nTime: %{time_total}s\\n" https://github.com \
                        && echo "✅ HTTPS 连接正常" \
                        || echo "❌ HTTPS 连接失败"

                    echo ""
                    echo "========== 环境变量（代理） =========="
                    echo "http_proxy=${http_proxy:-未设置}"
                    echo "https_proxy=${https_proxy:-未设置}"
                    echo "no_proxy=${no_proxy:-未设置}"
                '''
            }
        }

        stage('拉取代码 (Checkout)') {
            steps {
                echo "正在从 ${GIT_URL} 的 ${GIT_BRANCH} 分支拉取代码..."
                retry(3) {
                    checkout([
                        $class: 'GitSCM',
                        branches: [[name: "*/${GIT_BRANCH}"]],
                        doGenerateSubmoduleConfigurations: false,
                        extensions: [
                            [$class: 'CloneOption', shallow: true, depth: 1, noTags: true, honorRefspec: true, timeout: 20],
                            [$class: 'CheckoutOption', timeout: 20],
                            [$class: 'PruneStaleBranch']
                        ],
                        submoduleCfg: [],
                        userRemoteConfigs: [[
                            url: "${GIT_URL}",
                            refspec: "+refs/heads/${GIT_BRANCH}:refs/remotes/origin/${GIT_BRANCH}"
                        ]]
                    ])
                }
            }
        }

        stage('初始化版本号 (Init Version)') {
            steps {
                script {
                    env.GIT_SHORT_HASH = sh(script: 'git rev-parse --short HEAD', returnStdout: true).trim()
                    env.IMAGE_TAG = "v-${env.GIT_SHORT_HASH}"

                    env.IMAGE_BACKEND_BASE = "${env.HARBOR_REGISTRY}/${env.HARBOR_PROJECT}/server"
                    env.IMAGE_FRONTEND_BASE = "${env.HARBOR_REGISTRY}/${env.HARBOR_PROJECT}/web"

                    env.FULL_IMAGE_BACKEND = "${env.IMAGE_BACKEND_BASE}:${env.IMAGE_TAG}"
                    env.FULL_IMAGE_BACKEND_LATEST = "${env.IMAGE_BACKEND_BASE}:latest"

                    env.FULL_IMAGE_FRONTEND = "${env.IMAGE_FRONTEND_BASE}:${env.IMAGE_TAG}"
                    env.FULL_IMAGE_FRONTEND_LATEST = "${env.IMAGE_FRONTEND_BASE}:latest"

                    echo "🌟 本次统一构建版本号为: ${env.IMAGE_TAG}"
                }
            }
        }

        stage('登录 Harbor') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: "${HARBOR_CREDENTIAL_ID}",
                    usernameVariable: 'HARBOR_USER',
                    passwordVariable: 'HARBOR_PASS'
                )]) {
                    sh '''
                        set -e
                        echo "$HARBOR_PASS" | docker login ${HARBOR_REGISTRY} -u "$HARBOR_USER" --password-stdin
                    '''
                }
            }
        }



        stage('构建并推送镜像 (Build & Push)') {
            parallel {
                stage('🚀 后端 (Server)') {
                    steps {
                        echo "开始构建后端镜像: ${FULL_IMAGE_BACKEND}"
                        retry(3) {
                            sh '''
                                set -e
                                docker build --network host \
                                  -f server/Dockerfile \
                                  -t ${FULL_IMAGE_BACKEND} \
                                  -t ${FULL_IMAGE_BACKEND_LATEST} \
                                  server
                            '''
                        }
                        sh '''
                            set -e
                            docker push ${FULL_IMAGE_BACKEND}
                            docker push ${FULL_IMAGE_BACKEND_LATEST}
                        '''
                    }
                }

                stage('🎨 前端 (Web)') {
                    steps {
                        echo "开始构建前端镜像: ${FULL_IMAGE_FRONTEND}"
                        retry(3) {
                            sh '''
                                set -e
                                docker build --network host \
                                  -f web/Dockerfile \
                                  -t ${FULL_IMAGE_FRONTEND} \
                                  -t ${FULL_IMAGE_FRONTEND_LATEST} \
                                  web
                            '''
                        }
                        sh '''
                            set -e
                            docker push ${FULL_IMAGE_FRONTEND}
                            docker push ${FULL_IMAGE_FRONTEND_LATEST}
                        '''
                    }
                }
            }
        }
    }

    post {
        success {
            echo "🎉 前后端镜像（${env.IMAGE_TAG}）已成功构建并推送。"
        }
        failure {
            echo "❌ 流水线执行失败，请检查上方日志定位问题。"
        }
        always {
            sh 'docker logout ${HARBOR_REGISTRY} || true'
        }
    }
}
