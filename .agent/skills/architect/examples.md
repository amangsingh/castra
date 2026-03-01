# `castra` CLI Examples for the Architect Role

## Creating a New Project
`castra project add --role architect --name "New Mobile App" --description "A revolutionary app for task management."`

## Adding a Sprint to a Project
`castra sprint add --role architect --project 1 --name "Q3 Development" --start "2026-07-01" --end "2026-09-30"`

## Adding a Task to a Sprint
`castra task add --role architect --project 1 --sprint 1 --title "Design User Authentication Flow" --desc "Create database schema and API endpoints for user login and registration."`
