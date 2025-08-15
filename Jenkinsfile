pipeline {
  agent any
  options { timeout(time: 15, unit: 'MINUTES'); timestamps() /* ansiColor('xterm') 可选 */ }
  tools { jdk 'jdk17' }
  environment { APP_ENV = 'dev'; BUILD_ID = "${env.BUILD_NUMBER}" }
  parameters {
    choice(name: 'ENV', choices: ['dev','test','prod'], description: 'Target env')
    booleanParam(name: 'DEPLOY', defaultValue: false)
  }
  triggers { cron('H/30 * * * *') }

  stages {
    stage('Hello / Smoke') {
      steps {
        echo "Job: ${env.JOB_NAME}  Build: #${env.BUILD_NUMBER}  Node: ${env.NODE_NAME}"
        sh '''
          set -eux
          echo "Params => DEPLOY=${DEPLOY}  ENV=${ENV}"
          echo "Whoami / Uname:"
          whoami || true
          uname -a || true
          echo "--- Top 40 ENV ---"
          env | sort | sed -n '1,40p'
        '''
      }
    }

    stage('Checkout') {
      steps {
        script { currentBuild.displayName = "#${env.BUILD_NUMBER} ${params.ENV}" }
        // 即使拉代码失败，也让流水线继续，保证后续还能打印东西
        catchError(buildResult: 'UNSTABLE', stageResult: 'FAILURE') {
          checkout([$class: 'GitSCM',
            branches: [[name: '*/main']],         // 如用 master 改成 */master
            userRemoteConfigs: [[
              url: 'https://github.com/aaronwork1205/c--fft-microservice'
              // , credentialsId: 'your-github-cred-id' // 私有仓库时解除注释
            ]],
            extensions: [[$class: 'CleanBeforeCheckout'], [$class: 'PruneStaleBranch']]
          ])
        }
      }
    }

    stage('Build') {
      steps {
        sh '''
          set -eux
          mvn -v
          mvn -B package
        '''
      }
    }

    stage('Test') {
      steps {
        echo 'Running JUnit...'
        junit 'target/surefire-reports/*.xml'
      }
    }

    stage('Package') {
      steps {
        sh '''
          set -eux
          echo "Artifacts under target/:"
          ls -al target || true
        '''
        archiveArtifacts artifacts: 'target/*.jar', fingerprint: true
      }
    }

    stage('Approval')  {
      when { expression { params.DEPLOY && params.ENV == 'prod' } }
      steps { input message: 'Approve deploy to PROD?', ok: 'Deploy' }
    }

    stage('Deploy') {
      when { expression { params.DEPLOY } }
      steps {
        sh '''
          set -eux
          echo "deploy to ${ENV}"
        '''
      }
    }

    stage('Parallel')  {
      parallel {
        stage('Lint') { steps { sh 'set -x; echo run linter && sleep 1' } }
        stage('Unit') { steps { sh 'set -x; echo run unit tests && sleep 1' } }
      }
    }

    stage('Matrix') {
      matrix {
        axes { axis { name 'PY'; values '3.9','3.11' } }
        stages {
          stage('Run') {
            steps { sh 'set -x; echo "python $PY"' }
          }
        }
      }
    }
  }

  post {
    success { echo "OK #${env.BUILD_NUMBER} ✅" }
    failure { echo "FAILED #${env.BUILD_NUMBER} ❌" }
    always  {
      echo "Build URL: ${env.BUILD_URL}"
      sh '''
        set -eux
        echo "Workspace snapshot:"
        pwd
        ls -al
      '''
      cleanWs()
    }
  }
}
