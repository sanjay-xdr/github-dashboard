# Github Dashboard 
## visualize your repositry data on various parameters
 - Pull Requests
 - Workflow
 - Repo Stats


![image](https://github.com/user-attachments/assets/419b7f30-0325-4882-9f0e-44fc21e75b4a)

![image](https://github.com/user-attachments/assets/e40a8ad4-6cf1-4308-b31c-3baa3a7789be)

# Setting Up the Project Locally

Follow these steps to set up and run the project locally on your machine.

## Clone the Repository

First, clone the repository to your local machine using the following command:

```bash
git clone https://github.com/sanjay-xdr/github-dashboard.git
```

## Running the Backend

1. Navigate to the backend directory:

    ```bash
    cd backend
    ```

2. Navigate to the `cmd/web` directory:

    ```bash
    cd cmd/web
    ```

3. Run the backend application:

    ```bash
    go run .
    ```

## Running the Frontend

1. Navigate to the frontend directory (if not already in the root directory):

    ```bash
    cd frontend/keploy-dashboard
    ```
2. Install all the dependencies
   
     ```bash
    npm  i 
    ```  

3. Start the frontend application:

    ```bash
    npm run dev
    ```

You should now have both the backend and frontend running locally. Open your browser and navigate to the appropriate URL to view the application.
