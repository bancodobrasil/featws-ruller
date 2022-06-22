[![Go Reference](https://pkg.go.dev/badge/github.com/abu-lang/goabu.svg)](https://pkg.go.dev/github.com/bancodobrasil/featws-ruller)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/bancodobrasil/featws-ruller/blob/develop/LICENSE)
# **Featws Ruller** [![About_en](https://github.com/yammadev/flag-icons/blob/master/png/US.png?raw=true)](https://github.com/bancodobrasil/featws-transpiler/blob/develop/README.md)
- The featws-ruller project is an implmentation of the [grule-rule-engine](https://github.com/hyperjumptech/grule-rule-engine), used to evalute rulesheets(.grl) 

## Required Software
- You must have the **Go Programming Language** installed in your machine to run this project. You can get the official download [here](https://go.dev/doc/install).

- Clone the **featws-transpiler** repository to your local machine and make sure both transpiler and ruller project are in the same folder. You can find the featws-transpiler repository [here](https://github.com/bancodobrasil/featws-transpiler).

## Initializate the project
- Clone this project to your local machine.
- On _main.go_ folder (*../featws-ruller/main.go*), open your local terminal and type the command `go run main.go`, if your OS is windows, you can by build and run the executable `go build && ./featws-ruler.exe`, or if your OS is mac, type the command `go build -o ruller && ./ruller $@`.

## Testing diferents rulesheets
- Check if you have in your workspace the **featws-transpiler** and copy the path from .grl file for the new case, you can find that on the cases _tests_ -> cases
- Now just replace the env variable "FEATWS_RULLER_DEFAULT_RULES" on .env file on ruller, with the new path, and run like the instructions above.

## Testing rulesheet with resolvers
- To test if the resolver are loaded, you have to set the **featws-resolver-bridge** URL, on the .env file to.

## Load a rulesheet from remote source
- To load a rulesheet from a remote soure, just change the .env variable "FEATWS_RULLER_RESOURCE_LOADER_URL" pointed to your URL.


# Using main endpoints
_By default the port will be :8000_
- GET **http://localhost:YOURSETTEDPORT/** 
  - Will return  simple message to client such as: "FeatWS Ruller Works!!!"

- POST **http://localhost:YOURSEETEDPORT/api/v1/eval** 
  - On this end point you must have to pass a body, witch is the parameters setted by rulesheet folder on the featws-transpiler. Using case 0001 for example the body should be:
    ```json
       {
          "mynumber": "45"
       }
    ```
    - Sending the request, response should be something like that: (this response difined by .featws file on ruleshet folder, in that case is false because the condition is mynumber > 12)
    ```json 
        {
          "myboolfeat": false
        }
    ```
- GET **http://localhost:YOURSEETEDPORT/swagger/index.html**
  - On your browser, you can see the swagger documentation of the api.

- GET **http://localhost:YOURSEETEDPORT/health/live?full=1** 
  - This endpoint will check the live status of the application, just like that:
    ```json
    {
      "goroutine-threshold": "OK"
    }
    ```

- GET **http://localhost:YOURSEETEDPORT/health/ready?full=1**
  - This endpoint will check the ready status of the services used by ruller project
  ```json
  {
  "goroutine-threshold": "OK",
  "resolver-bridge": "OK",
  "resource-loader": "OK"
  }
  ```





      









