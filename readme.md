This directory contains information to build a docker repository, that can be
used to develop web applications using go and webpack.

The directory to use should be mounted as volume to the container and should
use the following structure:

- _build				Build directory containing all executable content
  - server				Directory for server binary and configuration files
  	- config.yaml		Application configuration in yaml format
  	- ... 				Executable files for the application
  - web					Directory for all public web content
    - ... 				Public content (html, javascript, css, assets)
  - Dockerfile			Definitions to create a docker container of the application
- server				Development directory for server content
  - vendor				Directory for external packages (installed with govendor)
  - .. 					Golang development files
- web					Development directory for web content
  - app					Frontend application code
  - node_modules		External packages for frontend code
  - package.json 		Specification of external packages
  - webpack.dev.config.js 		Configuration for webpack builder for development
  - webpack.build.config.js 	Configuration for webpack builder for production
- readme.md 			Description of the application and install and run instructions

Sample command to start the container
> docker run --rm -ti -v "$PWD/test-project:/app" -e "GOPACKAGE=bitbucket.com/dkfbasel/test-project" development
