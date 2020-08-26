## Learnings

Main - initialise app
     - run app

Within app.go.
Initialise the routes. Routes are registered to a handler, when a handler
matches on a given path a handler function is called. 
This function can take variables from the path. These app methods then invoke
the underlying model methods used to communicate to the database. 
