type OpenAPIObject = {
  openapi: string;
  info: {
    title: string;
    version: string;
    description?: string;
  };
  servers?: Array<{
    url: string;
    description?: string;
  }>;
  paths: Record<string, any>;
  components?: Record<string, any>;
};

export const swaggerDoc: OpenAPIObject = {
  openapi: "3.0.0",
  info: {
    title: "Training Records API",
    version: "1.0.0",
    description: "API for managing training records, exercises, and sets",
  },
  servers: [
    {
      url: "https://training-api.workers.dev",
      description: "Production server",
    },
    {
      url: "http://localhost:8787",
      description: "Local development server",
    },
  ],
  paths: {
    "/auth/register": {
      post: {
        summary: "Register a new user",
        tags: ["Authentication"],
        requestBody: {
          required: true,
          content: {
            "application/json": {
              schema: {
                type: "object",
                required: ["email", "password", "name"],
                properties: {
                  email: {
                    type: "string",
                    format: "email",
                    description: "User's email address",
                  },
                  password: {
                    type: "string",
                    format: "password",
                    minLength: 8,
                    description: "User's password (min 8 characters)",
                  },
                  name: {
                    type: "string",
                    description: "User's full name",
                  },
                },
              },
            },
          },
        },
        responses: {
          "201": {
            description: "User registered successfully",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    message: {
                      type: "string",
                      example: "User registered successfully",
                    },
                    token: {
                      type: "string",
                      example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
                    },
                    user: {
                      type: "object",
                      properties: {
                        id: {
                          type: "string",
                          format: "uuid",
                        },
                        email: {
                          type: "string",
                          format: "email",
                        },
                        name: {
                          type: "string",
                        },
                      },
                    },
                  },
                },
              },
            },
          },
          "409": {
            description: "User with this email already exists",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "User with this email already exists",
                    },
                  },
                },
              },
            },
          },
          "500": {
            description: "Server error",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "Failed to register user",
                    },
                  },
                },
              },
            },
          },
        },
      },
    },
    "/auth/login": {
      post: {
        summary: "Login with email and password",
        tags: ["Authentication"],
        requestBody: {
          required: true,
          content: {
            "application/json": {
              schema: {
                type: "object",
                required: ["email", "password"],
                properties: {
                  email: {
                    type: "string",
                    format: "email",
                    description: "User's email address",
                  },
                  password: {
                    type: "string",
                    format: "password",
                    description: "User's password",
                  },
                },
              },
            },
          },
        },
        responses: {
          "200": {
            description: "Login successful",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    message: {
                      type: "string",
                      example: "Login successful",
                    },
                    token: {
                      type: "string",
                      example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
                    },
                    user: {
                      type: "object",
                      properties: {
                        id: {
                          type: "string",
                          format: "uuid",
                        },
                        email: {
                          type: "string",
                          format: "email",
                        },
                        name: {
                          type: "string",
                        },
                      },
                    },
                  },
                },
              },
            },
          },
          "401": {
            description: "Invalid credentials",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "Invalid credentials",
                    },
                  },
                },
              },
            },
          },
          "500": {
            description: "Server error",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "Failed to login",
                    },
                  },
                },
              },
            },
          },
        },
      },
    },
    "/auth/me": {
      get: {
        summary: "Get current user information",
        tags: ["Authentication"],
        security: [
          {
            bearerAuth: [],
          },
        ],
        responses: {
          "200": {
            description: "User information retrieved successfully",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    user: {
                      type: "object",
                      properties: {
                        id: {
                          type: "string",
                          format: "uuid",
                        },
                        email: {
                          type: "string",
                          format: "email",
                        },
                        name: {
                          type: "string",
                        },
                        created_at: {
                          type: "string",
                          format: "date-time",
                        },
                      },
                    },
                  },
                },
              },
            },
          },
          "401": {
            description: "Unauthorized",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "Unauthorized",
                    },
                  },
                },
              },
            },
          },
          "404": {
            description: "User not found",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "User not found",
                    },
                  },
                },
              },
            },
          },
        },
      },
    },
    "/api/v1/records": {
      get: {
        summary: "Get all training records",
        tags: ["Training Records"],
        security: [
          {
            bearerAuth: [],
          },
        ],
        responses: {
          "200": {
            description: "List of training records",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    records: {
                      type: "array",
                      items: {
                        $ref: "#/components/schemas/TrainingRecord",
                      },
                    },
                  },
                },
              },
            },
          },
          "500": {
            description: "Server error",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "Failed to fetch training records",
                    },
                  },
                },
              },
            },
          },
        },
      },
      post: {
        summary: "Create a new training record",
        tags: ["Training Records"],
        security: [
          {
            bearerAuth: [],
          },
        ],
        requestBody: {
          required: true,
          content: {
            "application/json": {
              schema: {
                type: "object",
                required: ["title", "date"],
                properties: {
                  title: {
                    type: "string",
                    description: "Title of the training record",
                  },
                  date: {
                    type: "string",
                    format: "date",
                    description: "Date of the training session",
                  },
                  description: {
                    type: "string",
                    description: "Description of the training session",
                  },
                  exercises: {
                    type: "array",
                    items: {
                      type: "object",
                      required: ["name"],
                      properties: {
                        name: {
                          type: "string",
                          description: "Name of the exercise",
                        },
                        sets: {
                          type: "array",
                          items: {
                            type: "object",
                            required: ["weight", "reps"],
                            properties: {
                              weight: {
                                type: "number",
                                description: "Weight used for the set",
                              },
                              reps: {
                                type: "integer",
                                description: "Number of repetitions",
                              },
                              notes: {
                                type: "string",
                                description: "Notes about the set",
                              },
                            },
                          },
                        },
                      },
                    },
                  },
                },
              },
            },
          },
        },
        responses: {
          "201": {
            description: "Training record created successfully",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    id: {
                      type: "string",
                      format: "uuid",
                    },
                    message: {
                      type: "string",
                      example: "Training record created successfully",
                    },
                  },
                },
              },
            },
          },
          "400": {
            description: "Invalid input",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "Title and date are required",
                    },
                  },
                },
              },
            },
          },
          "500": {
            description: "Server error",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "Failed to create training record",
                    },
                  },
                },
              },
            },
          },
        },
      },
    },
    "/api/v1/records/{id}": {
      get: {
        summary: "Get a specific training record",
        tags: ["Training Records"],
        security: [
          {
            bearerAuth: [],
          },
        ],
        parameters: [
          {
            name: "id",
            in: "path",
            required: true,
            schema: {
              type: "string",
              format: "uuid",
            },
            description: "ID of the training record",
          },
        ],
        responses: {
          "200": {
            description: "Training record details",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    record: {
                      $ref: "#/components/schemas/TrainingRecordDetail",
                    },
                  },
                },
              },
            },
          },
          "404": {
            description: "Training record not found",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "Training record not found",
                    },
                  },
                },
              },
            },
          },
          "500": {
            description: "Server error",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "Failed to fetch training record",
                    },
                  },
                },
              },
            },
          },
        },
      },
      put: {
        summary: "Update a training record",
        tags: ["Training Records"],
        security: [
          {
            bearerAuth: [],
          },
        ],
        parameters: [
          {
            name: "id",
            in: "path",
            required: true,
            schema: {
              type: "string",
              format: "uuid",
            },
            description: "ID of the training record",
          },
        ],
        requestBody: {
          required: true,
          content: {
            "application/json": {
              schema: {
                type: "object",
                required: ["title", "date"],
                properties: {
                  title: {
                    type: "string",
                    description: "Title of the training record",
                  },
                  date: {
                    type: "string",
                    format: "date",
                    description: "Date of the training session",
                  },
                  description: {
                    type: "string",
                    description: "Description of the training session",
                  },
                },
              },
            },
          },
        },
        responses: {
          "200": {
            description: "Training record updated successfully",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    message: {
                      type: "string",
                      example: "Training record updated successfully",
                    },
                  },
                },
              },
            },
          },
          "400": {
            description: "Invalid input",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "Title and date are required",
                    },
                  },
                },
              },
            },
          },
          "404": {
            description: "Training record not found",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "Training record not found",
                    },
                  },
                },
              },
            },
          },
          "500": {
            description: "Server error",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "Failed to update training record",
                    },
                  },
                },
              },
            },
          },
        },
      },
      delete: {
        summary: "Delete a training record",
        tags: ["Training Records"],
        security: [
          {
            bearerAuth: [],
          },
        ],
        parameters: [
          {
            name: "id",
            in: "path",
            required: true,
            schema: {
              type: "string",
              format: "uuid",
            },
            description: "ID of the training record",
          },
        ],
        responses: {
          "200": {
            description: "Training record deleted successfully",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    message: {
                      type: "string",
                      example: "Training record deleted successfully",
                    },
                  },
                },
              },
            },
          },
          "404": {
            description: "Training record not found",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "Training record not found",
                    },
                  },
                },
              },
            },
          },
          "500": {
            description: "Server error",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "Failed to delete training record",
                    },
                  },
                },
              },
            },
          },
        },
      },
    },
    "/api/v1/records/{id}/exercises": {
      post: {
        summary: "Add an exercise to a training record",
        tags: ["Exercises"],
        security: [
          {
            bearerAuth: [],
          },
        ],
        parameters: [
          {
            name: "id",
            in: "path",
            required: true,
            schema: {
              type: "string",
              format: "uuid",
            },
            description: "ID of the training record",
          },
        ],
        requestBody: {
          required: true,
          content: {
            "application/json": {
              schema: {
                type: "object",
                required: ["name"],
                properties: {
                  name: {
                    type: "string",
                    description: "Name of the exercise",
                  },
                  sets: {
                    type: "array",
                    items: {
                      type: "object",
                      required: ["weight", "reps"],
                      properties: {
                        weight: {
                          type: "number",
                          description: "Weight used for the set",
                        },
                        reps: {
                          type: "integer",
                          description: "Number of repetitions",
                        },
                        notes: {
                          type: "string",
                          description: "Notes about the set",
                        },
                      },
                    },
                  },
                },
              },
            },
          },
        },
        responses: {
          "201": {
            description: "Exercise added successfully",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    id: {
                      type: "string",
                      format: "uuid",
                    },
                    message: {
                      type: "string",
                      example: "Exercise added successfully",
                    },
                  },
                },
              },
            },
          },
          "400": {
            description: "Invalid input",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "Exercise name is required",
                    },
                  },
                },
              },
            },
          },
          "404": {
            description: "Training record not found",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "Training record not found",
                    },
                  },
                },
              },
            },
          },
          "500": {
            description: "Server error",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "Failed to add exercise",
                    },
                  },
                },
              },
            },
          },
        },
      },
    },
    "/api/v1/exercises/{id}/sets": {
      post: {
        summary: "Add a set to an exercise",
        tags: ["Sets"],
        security: [
          {
            bearerAuth: [],
          },
        ],
        parameters: [
          {
            name: "id",
            in: "path",
            required: true,
            schema: {
              type: "string",
              format: "uuid",
            },
            description: "ID of the exercise",
          },
        ],
        requestBody: {
          required: true,
          content: {
            "application/json": {
              schema: {
                type: "object",
                required: ["weight", "reps"],
                properties: {
                  weight: {
                    type: "number",
                    description: "Weight used for the set",
                  },
                  reps: {
                    type: "integer",
                    description: "Number of repetitions",
                  },
                  notes: {
                    type: "string",
                    description: "Notes about the set",
                  },
                },
              },
            },
          },
        },
        responses: {
          "201": {
            description: "Set added successfully",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    id: {
                      type: "string",
                      format: "uuid",
                    },
                    message: {
                      type: "string",
                      example: "Set added successfully",
                    },
                  },
                },
              },
            },
          },
          "400": {
            description: "Invalid input",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "Weight and reps are required",
                    },
                  },
                },
              },
            },
          },
          "404": {
            description: "Exercise not found",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "Exercise not found",
                    },
                  },
                },
              },
            },
          },
          "500": {
            description: "Server error",
            content: {
              "application/json": {
                schema: {
                  type: "object",
                  properties: {
                    error: {
                      type: "string",
                      example: "Failed to add set",
                    },
                  },
                },
              },
            },
          },
        },
      },
    },
  },
  components: {
    securitySchemes: {
      bearerAuth: {
        type: "http",
        scheme: "bearer",
        bearerFormat: "JWT",
      },
    },
    schemas: {
      TrainingRecord: {
        type: "object",
        properties: {
          id: {
            type: "string",
            format: "uuid",
          },
          title: {
            type: "string",
          },
          date: {
            type: "string",
            format: "date",
          },
          description: {
            type: "string",
          },
          created_at: {
            type: "string",
            format: "date-time",
          },
          updated_at: {
            type: "string",
            format: "date-time",
          },
          user_id: {
            type: "string",
            format: "uuid",
          },
        },
      },
      TrainingRecordDetail: {
        type: "object",
        properties: {
          id: {
            type: "string",
            format: "uuid",
          },
          title: {
            type: "string",
          },
          date: {
            type: "string",
            format: "date",
          },
          description: {
            type: "string",
          },
          created_at: {
            type: "string",
            format: "date-time",
          },
          updated_at: {
            type: "string",
            format: "date-time",
          },
          user_id: {
            type: "string",
            format: "uuid",
          },
          exercises: {
            type: "array",
            items: {
              type: "object",
              properties: {
                id: {
                  type: "string",
                  format: "uuid",
                },
                record_id: {
                  type: "string",
                  format: "uuid",
                },
                name: {
                  type: "string",
                },
                sets: {
                  type: "array",
                  items: {
                    type: "object",
                    properties: {
                      id: {
                        type: "string",
                        format: "uuid",
                      },
                      exercise_id: {
                        type: "string",
                        format: "uuid",
                      },
                      weight: {
                        type: "number",
                      },
                      reps: {
                        type: "integer",
                      },
                      notes: {
                        type: "string",
                      },
                    },
                  },
                },
              },
            },
          },
        },
      },
    },
  },
};
