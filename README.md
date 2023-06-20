
# My Gin-MVC-CLI TUI Project

This is a command-line tool for creating Gin projects structured in a pattern similar to the Model-View-Controller (MVC) pattern using Bubbletea. It provides an interactive TUI (Text User Interface) for generating the project folder structure and files.

## Installation

To install and use the tool, follow these steps:

1. Clone the repository:
   ```
   git clone https://github.com/steelthedev/gin-mvc-cli.git
   ```

2. Change into the project directory:
   ```
   cd gin-mvc-cli
   ```

3. Build and install the application using `go install`:
   ```
   go install .
   ```

## Usage

To start the TUI and create a new Gin-MVC project, run the following command:
```
gin-mvc-cli
```

The TUI will guide you through the process of setting up the project structure, including creating the necessary folders and files.

## Project Structure

The project follows a structured directory layout to separate concerns and maintain a clean codebase. Here is an overview of the generated project structure:

- `controllers/`: Contains the controller logic for handling requests and generating responses.
- `models/`: Includes the data models used in the application.
- `routes/`: Contains all the neccessary routes.
- `main.go`: The entry point of the application.

Feel free to modify and expand the structure according to your project's needs.

## Dependencies

The project relies on the following dependencies:

- `github.com/gin-gonic/gin`: The Gin web framework for handling HTTP requests.
- `https://github.com/charmbracelet/bubbletea.git`: The bubbletea project
Make sure you have this dependency installed and properly configured in your environment.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, feel free to open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
```


[![Demo Video](https://drive.google.com/file/d/1kucMWvTEbB9ej0ZeQmnANQP8PZxNOptN/view?usp=drive_link)](https://drive.google.com/file/d/1kucMWvTEbB9ej0ZeQmnANQP8PZxNOptN/view?usp=drive_link)
