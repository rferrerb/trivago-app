node {

    checkout scm

    env.DOCKER_API_VERSION="1.23"

    sh "git rev-parse --short HEAD > commit-id"

    tag = readFile('commit-id').replace("\n", "").replace("\r", "")
    appName = "trivago-app"
    registryHost = "127.0.0.1:30400/"
    imageName = "${registryHost}${appName}:${tag}"
    env.BUILDIMG=imageName

    stage "Build"
        sh 'go get -d -v'
        sh 'GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main'
        sh "docker build -t ${imageName} -f Dockerfile ."

    stage "Push"

        sh "docker push ${imageName}"

    stage "Deploy"

        kubernetesDeploy configs: "${appName}/k8s/*.yaml", kubeconfigId: 'k8s-local'

}