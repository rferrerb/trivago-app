node {

    checkout scm

    env.DOCKER_API_VERSION="1.23"

    sh "git rev-parse --short HEAD > commit-id"

    tag = readFile('commit-id').replace("\n", "").replace("\r", "")
    appName = "trivago-app"
    registryHost = "127.0.0.1:30400/"
    imageName = "${registryHost}${appName}:${tag}"
    env.BUILDIMG=imageName
    withCredentials([usernamePassword(credentialsId: 'abb9982a-446d-4917-95f7-85bb98142d18', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD')])
    stage "Build"
        sh 'go get -d -v'
        sh 'CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .'
        sh 'docker build --build-arg database_user=$USERNAME --build-arg database_password=$PASSWORD   -t ${imageName} -f Dockerfile .'

    stage "Push"

        sh "docker push ${imageName}"

    stage "Deploy"
        sh "sed -i 's/\$BUILD_NUMBER/${tag}/' ${appName}/k8s/deployment.yaml "
        kubernetesDeploy configs: "${appName}/k8s/*.yaml", kubeconfigId: 'k8s-local'

}