# Procurement System

This project is a web-based procurement system designed to manage requisitions, tenders, purchase orders, and other related workflows.

## Technology Stack

- **Frontend:** SvelteKit, TypeScript, Tailwind CSS
- **Backend:** Go (planned)
- **Database:** SQLite (planned)

## Current Status

- Basic project structure set up.
- Initial UI for layout, navigation, and stubs for core modules (requisitions, tenders) implemented.
- Placeholder icons are in use due to temporary issues with an icon library.

## Next Steps

- Continue development of core modules.
- Integrate with backend services.
- Implement user authentication.
- Refine UI and add more features as per the project checklist.

## Getting Started

This section will guide you through setting up and running the Procurement System application locally.

### Prerequisites

Ensure you have the following installed on your system:

- [Node.js](https://nodejs.org/) (which includes npm) for the frontend.
- [Go](https://golang.org/doc/install) for the backend.

### Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/gTahidi/Procure.git
    cd procurement
    ```

2.  **Set up the Frontend:**
    Navigate to the `frontend` directory and install the necessary Node.js dependencies:
    ```bash
    cd frontend
    npm install
    cd ..
    ```

3.  **Set up the Backend:**
    Navigate to the `backend` directory and ensure Go modules are tidy (this might download dependencies defined in `go.mod`):
    ```bash
    cd backend
    go mod tidy
    cd ..
    ```

## Running the Application

To run the application, you'll need to start both the backend and frontend services.

### 1. Initialize the Database (Backend)

If this is your first time running the application or if you need to reset the database:

1.  Navigate to the `backend` directory:
    ```bash
    cd backend
    ```
2.  Run the database initialization script:
    ```bash
    go run db_init.go
    ```
3.  Navigate back to the project root:
    ```bash
    cd ..
    ```
    This command will execute the `db_init.go` file, which should set up the `procurement.db` SQLite database in the `backend` directory.

### 2. Start the Backend Server

1.  Navigate to the `backend` directory:
    ```bash
    cd backend
    ```
2.  **[TODO: Add command to start the Go backend server here]**
    For example, if you have a `main.go` that starts the server, it might be `go run main.go` or simply `./your-compiled-backend-app`.
3.  The backend server should now be running (typically on a port like `localhost:8080` or similar - please specify).

### 3. Start the Frontend Server

1.  In a new terminal window or tab, navigate to the `frontend` directory:
    ```bash
    cd frontend
    ```
2.  Start the SvelteKit development server:
    ```bash
    npm run dev
    ```
3.  The frontend application will typically be available at `http://localhost:5173` (or another port specified by SvelteKit).

Once both the backend and frontend servers are running, you should be able to access the Procurement System application in your web browser via the frontend's URL.
