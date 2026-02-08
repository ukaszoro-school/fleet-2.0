# Fleet manager
A web app for managing and visualising bus routes. Simple menu for making and removing new stops and grouping them together with times of arrival to create simple routes.

### Building
```
cd server
go build
```

### Running
Use the created executable created in the server directory.

## Technology
- Backend: Golang
- Frontend: htmx
- Database: MongoDB

I chose Go, MongoDB, and HTMX to keep the application simple, fast, and maintainable while avoiding unnecessary frontend complexity. Go provides strong performance, type safety, and a straightforward standard library that works well for building reliable HTTP services. MongoDB fits naturally with the application’s data model, allowing flexible schemas for routes and stops without excessive migration overhead. HTMX enables dynamic, responsive user interactions directly from server-rendered HTML, eliminating the need for a heavy JavaScript framework while keeping the frontend tightly coupled to backend logic. Together, this stack prioritizes clarity, server-driven behavior, and long-term maintainability.

## ERD
![alt text](https://github.com/ukaszoro-school/fleet-2.0/blob/main/db.png)
