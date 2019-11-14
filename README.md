# ApiBoy Backend

## Introduction

ApiBoy is an application for testing APIs, very similar to other tools like Postman and Insomnia, but with some important differences:

- ApiBoy is 100% open source.
- ApiBoy is a serverless web application with a simplified set of features in comparison with other tools.
- ApiBoy can be easily deployed to your own infrastructure at no cost because it uses cloud services that have a comfortable free tier: AWS API Gateway, AWS Lambda, Firebase Realtime Database, Firebase Firestore and Netlify.
- ApiBoy has real time features for teams to work and collaborate in the same projects at the same time.

The project is composed by a Backend (this repository) and a [Frontend](https://github.com/api-boy/frontend).

The Backend is written in [Go](https://golang.org/) and uses [Firebase Firestore](https://firebase.google.com/docs/firestore/) to store data.

The Frontend is implemented with [VueJS](https://vuejs.org/) and connects to Firebase in order to read data and receive changes in real time, while writes are done through the Backend by calling its endpoints.

## Demo

<img src="https://user-images.githubusercontent.com/8256604/68883480-e9578980-06ef-11ea-88b6-91fbee2a1336.gif" width="920" height="500">

## Setup

### Install `go`:

Follow the official installation instructions:
- [How to install go compiler](https://golang.org/doc/install)

### Install `up`:

Follow the official installation instructions:
- [How to install up](https://up.docs.apex.sh/)

### Install `aws` cli:

Follow the official installation instructions:
- [How to install aws cli](https://docs.aws.amazon.com/cli/latest/userguide/installing.html)

### Configure AWS credentials profile for the project:

Add the `apiboy` profile to the `~/.aws/credentials` file with your AWS access and secret keys.

```
[apiboy]
aws_access_key_id = xxxxxxxx
aws_secret_access_key = xxxxxxxxxxxxxxxxxxxxxxxx
```

### Install development tools:

```bash
make tools
```

### Install project dependencies:

```bash
make deps
```

### Configure project:

The project setup is done with [team](https://github.com/andybar2/team), that is installed with the development tools. **You only need to do these steps once**, and the rest of your team developers pointing to the same AWS account will automatically inherit the project configuration by just running or deploying the project.

Create a `.team` folder for each stage:

```bash
mkdir -p .team/development
mkdir -p .team/production
```

Create two [Firebase](https://firebase.google.com/) projects for `development` and `production` respectively. Generate a Firebase service account file for each of the projects. Then upload the Firebase service account files for each stage:

```bash
# Place the Firebase service account file for development
# in .team/development/firebase-service-account.json and
# then upload it with the following command:
team files upload -s "development" -p ".team/development/firebase-service-account.json"

# Place the Firebase service account file for production
# in .team/production/firebase-service-account.json and
# then upload it with the following command:
team files upload -s "production" -p ".team/production/firebase-service-account.json"
```

Set environment variables for each stage:

```bash
# Set log level:
team env set -s "development" -n "LOG_LEVEL" -v "debug"
team env set -s "production" -n "LOG_LEVEL" -v "info"

# Set firebase project id:
team env set -s "development" -n "FIREBASE_PROJECT_ID" -v "abiboy-dev-xxxxx"
team env set -s "production" -n "FIREBASE_PROJECT_ID" -v "abiboy-prod-zzzzz"

# Set firebase api key:
team env set -s "development" -n "FIREBASE_API_KEY" -v "XXXXXXXXXX"
team env set -s "production" -n "FIREBASE_API_KEY" -v "ZZZZZZZZZZ"

# Set jwt issuer:
team env set -s "development" -n "JWT_ISSUER" -v "abiboy-dev"
team env set -s "production" -n "JWT_ISSUER" -v "abiboy-prod"

# Set jwt sign key:
team env set -s "development" -n "JWT_SIGN_KEY" -v "XXXXXXXXXX"
team env set -s "production" -n "JWT_SIGN_KEY" -v "ZZZZZZZZZZ"
```

Configure the access rules for the _Firestore Database_ with the following code:

```
service cloud.firestore {
  match /databases/{database}/documents {
  	function signedIn() {
      return request.auth.uid != null;
    }
    
    function validProjectForUser(projectId) {
      return exists(/databases/$(database)/documents/projectusers/$(projectId + "-" + request.auth.uid));
    }
  
    match /{document=**} {
      allow read, write: if false;
    }
    
    match /users/{userId} {
      allow read: if signedIn() && resource.data.id == request.auth.uid;
    }
    
    match /projects/{projectId} {
      allow read: if signedIn() && validProjectForUser(resource.data.id);
    }
    
    match /projectusers/{projectUserId} {
      allow read: if signedIn() && resource.data.user_id == request.auth.uid;
    }
    
    match /folders/{folderId} {
      allow read: if signedIn() && validProjectForUser(resource.data.project_id);
    }
    
    match /requests/{requestId} {
      allow read: if signedIn() && validProjectForUser(resource.data.project_id);
    }

    match /environments/{environmentId} {
      allow read: if signedIn() && validProjectForUser(resource.data.project_id);
    }
  }
}
```

Configure the access rules for the _Realtime Database_ with the following code:

```
{
  "rules": {
    "presence": {
      ".read": "auth != null",
      ".write": "auth != null",
      "$projectId": {
        "$deviceId": {
          ".validate": "newData.hasChildren(['id', 'name', 'email'])",
          "id": {
            ".validate": "newData.val() == auth.uid"
          },
          "name": {
            ".validate": "newData.val() == auth.token.user_name"
          },
          "email": {
            ".validate": "newData.val() == auth.token.user_email"
          },
          "$other": {
            ".validate": false
          }
        }
      }
    }
  }
}
```

## Execution and Deployment

### Run development stage:

```bash
make dev
```

### Deploy production stage:

```bash
make production
```
