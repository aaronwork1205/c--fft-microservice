pipeline {
  agent any                       // 运行在哪：any / { label 'linux' } / { docker { image 'node:20' } }
  options { timeout(time: 15, unit: 'MINUTES'); timestamps() }
  tools { jdk 'jdk17' }           // 自动安装/声明工具
  environment { APP_ENV = 'dev'; BUILD_ID = "${env.BUILD_NUMBER}" }
  parameters {
    choice(name: 'ENV', choices: ['dev','test','prod'], description: 'Target env')
    booleanParam(name: 'DEPLOY', defaultValue: false)
  }
  triggers { cron('H/30 * * * *') }   // 可选：轮询/定时
  stages {
    stage('Checkout') { steps { checkout scm } }
    stage('Build')     { steps { sh 'mvn -B package' } }
    stage('Test')      { steps { junit 'target/surefire-reports/*.xml' } }
    stage('Package')   { steps { archiveArtifacts artifacts: 'target/*.jar', fingerprint: true } }
    stage('Approval')  { when { expression { params.DEPLOY && params.ENV == 'prod' } }
                         steps { input message: 'Approve deploy to PROD?', ok: 'Deploy' } }
    stage('Deploy')    { when { expression { params.DEPLOY } }
                         steps { sh 'echo deploy to ${ENV}' } }
    stage('Parallel')  {
      parallel {
        stage('Lint') { steps { sh 'echo run linter && sleep 1' } }
        stage('Unit') { steps { sh 'echo run unit tests && sleep 1' } }
      }
    }
    stage('Matrix') {
      matrix {
        axes { axis { name 'PY'; values '3.9','3.11' } }
        stages { stage('Run') { steps { sh 'echo "python $PY"' } } }
      }
    }
  }
  post {
    success { echo "OK #${env.BUILD_NUMBER}" }
    failure { echo "FAILED #${env.BUILD_NUMBER}" }
    always  { cleanWs() }
  }
}
