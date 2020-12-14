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


## Debugging in a container (VS Code)

When opening the folder, a popup should prompt to "Re-open in container" if it does not, then from the `F1` menu select "Remote-Containers: Re-open in container"

Open chrome with remote debugging:

```bash
google-chrome-stable --remote-debugging-port=9222 --user-data-dir=/tmp/remote-profile
```

Run the compound launch config (API + NG + Attach)