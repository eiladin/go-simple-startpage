# Go-Simple-Startpage

This is an example app built with go and angular 8.

## Building / Running 

To run this app, simply clone the repo, enter the `ui/` folder and run `npm install`

After the install is complete, you can open the root folder in VS Code, open main.go and hit F5 to start the API, alternatively, you can run `go run main.go` from the root folder to start the API.

In a terminal window, open the `ui/` folder and run `ng serve`

Navigate to http://localhost:4200 in your browser.

## Testing with go serving the angular app

In the `ui/` folder run `ng build --prod` (this will build the front-end bundle)

Once, that is complete, in the root folder run `go run main.go` or open the root folder in VS Code, open main.go and hit F5.