pipeline {
  agent any
  stages {
    stage('Hello') {
      steps {
        echo "Hi ${env.BUILD_USER ?: 'Jenkins'}!"
        sh '''
          echo "WORKSPACE=$WORKSPACE"
          echo "BRANCH_NAME=${BRANCH_NAME:-(no scm)}"
        '''
      }
    }
  }
  post {
    always { echo "Build #${env.BUILD_NUMBER} finished at ${new Date()}" }
  }
}
