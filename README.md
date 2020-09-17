# Go-Simple-Startpage

This is an example app built with go and angular 10.

## Building / Running 

To run this app, simply clone the repo, enter the `ui/` folder and run `npm install`

After the install is complete:
 - start angular with `ng serve` in the ui folder.
 - start the go api with `go run main.go` or `air` in the root folder
- Navigate to http://localhost:4200 in your browser.

## Testing with go serving the angular app

In the `ui/` folder run `ng build --prod` (this will build the front-end bundle)

Once, that is complete, in the root folder run `go run main.go`, `air`, or open the root folder in VS Code, open main.go and hit F5.
