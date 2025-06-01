# JoinTrip 🌍✈️

A social travel platform that connects travelers, facilitates trip sharing, and manages shared expenses.

## Overview

JoinTrip is a comprehensive travel companion application built in Go that allows users to:
- Find travel companions with similar destinations and interests
- Share and split travel expenses seamlessly
- Communicate with fellow travelers
- Discover and organize group trips
- Build trust through user ratings and reviews

## 🚀 Features

### Core Functionality
- **User Authentication & Profiles**: Secure registration, login, and comprehensive user profiles with verification
- **Trip Management**: Create, search, and manage travel opportunities with detailed filtering
- **Social Networking**: Connect with like-minded travelers and build a travel network
- **Expense Sharing**: Advanced expense splitting and tracking with multi-currency support
- **Communication**: Real-time messaging and trip-specific discussions
- **Rating System**: Build trust through peer reviews and reputation scores

### Advanced Features
- **Smart Search**: Find trips by destination, dates, budget, activities, and more
- **Geographic Integration**: Interactive maps and location-based services
- **Multi-Currency Support**: Handle expenses in different currencies with conversion
- **Privacy Controls**: Granular privacy settings for profiles and trips
- **Notification System**: Real-time updates for trip requests, messages, and activities
- **Mobile Responsive**: Optimized for mobile devices and progressive web app capabilities

## 🛠️ Technology Stack

### Backend
- **Language**: Go (Golang) 1.24+
- **Web Framework**: Gin/Echo (planned)
- **Database**: PostgreSQL with sqlx
- **Authentication**: Google OAuth 2.0 + JWT sessions
- **API**: RESTful API design
- **Static Assets**: Go embed for serving React build

### Frontend
- **Framework**: React 18+
- **Build Tool**: Vite/Create React App
- **State Management**: React Context/Redux Toolkit (planned)
- **Styling**: Tailwind CSS/Material-UI (planned)
- **HTTP Client**: Axios/Fetch API
- **Routing**: React Router

### Architecture
- **Deployment**: Single binary with embedded React build
- **Asset Serving**: Go `embed` directive for static files
- **API Communication**: JSON REST API between React and Go
- **Build Process**: React build embedded into Go binary

### Development Tools
- **Documentation**: Markdown with Mermaid diagrams
- **Version Control**: Git
- **Package Management**: Go modules + npm/yarn
- **Hot Reload**: Air (Go) + Vite/CRA dev server (React)

## 📁 Project Structure

```
jointrip/
├── docs/                              # Documentation
│   ├── entity_relationship_model.md     # Database design
│   ├── issues.md                        # Development issues tracking
│   └── sprint_plan.md                  # Sprint planning
├── cmd/                               # Application entry points
│   └── server/                        # Main server application
├── internal/                          # Private application code
│   ├── auth/                         # Authentication logic
│   ├── handlers/                     # HTTP handlers (API endpoints)
│   ├── models/                       # Data models
│   ├── services/                     # Business logic
│   ├── database/                     # Database operations
│   ├── middleware/                   # HTTP middleware
│   └── config/                       # Configuration management
├── pkg/                              # Public packages
├── migrations/                       # Database migrations
├── web/                              # Frontend React application
│   ├── public/                       # Public assets
│   ├── src/                          # React source code
│   │   ├── components/               # Reusable components
│   │   ├── pages/                    # Page components
│   │   ├── hooks/                    # Custom React hooks
│   │   ├── services/                 # API service calls
│   │   ├── utils/                    # Utility functions
│   │   ├── styles/                   # CSS/styling files
│   │   └── App.js                    # Main App component
│   ├── package.json                  # Node.js dependencies
│   ├── vite.config.js               # Vite configuration
│   └── dist/                         # Built React files (embedded)
├── static/                           # Additional static assets
├── tests/                            # Test files
│   ├── api/                          # API tests
│   └── integration/                  # Integration tests
├── scripts/                          # Build and deployment scripts
│   ├── build.sh                      # Build script
│   └── dev.sh                        # Development script
├── go.mod                            # Go module definition
├── go.sum                            # Go module checksums
├── main.go                           # Application entry point
├── embed.go                          # Go embed directives
└── README.md                         # This file
```

## 🏗️ Development Roadmap

The project is organized into 10 sprints, each lasting 2 weeks:

### Sprint 1: Setup and Authentication ✅
- [x] Project initialization
- [x] Basic Go structure
- [ ] React frontend setup with Vite/CRA
- [ ] Go embed configuration for React build
- [ ] Google OAuth 2.0 setup (Go backend)
- [ ] Database setup with PostgreSQL + sqlx
- [ ] JWT session management
- [ ] React Google OAuth integration
- [ ] Protected routes and auth middleware
- [ ] API client setup in React

### Sprint 2: User Profiles
- [ ] Profile API endpoints (Go backend)
- [ ] React profile components and forms
- [ ] Photo upload functionality (backend + frontend)
- [ ] Identity verification system
- [ ] Reputation/rating system
- [ ] Profile management UI in React

### Sprint 3: Trip Management
- [ ] Trip creation and editing
- [ ] Advanced search and filtering
- [ ] Geographic integration
- [ ] Tagging system

### Sprint 4: Communication
- [ ] Messaging system
- [ ] Real-time notifications
- [ ] Trip join requests
- [ ] Comment system

### Sprint 5: Expense Management
- [ ] Expense tracking
- [ ] Splitting calculator
- [ ] Multi-currency support
- [ ] Balance visualization

### Sprint 6: Security and Privacy
- [ ] Privacy controls
- [ ] Two-factor authentication
- [ ] User reporting system
- [ ] Security audit

### Sprint 7: External Integrations
- [ ] Travel service APIs
- [ ] Social media sharing
- [ ] Recommendation engine
- [ ] Weather and events

### Sprint 8: Testing and Optimization
- [ ] Comprehensive testing
- [ ] Performance optimization
- [ ] Security testing
- [ ] Bug fixes

### Sprint 9: Mobile Version
- [ ] Mobile optimization
- [ ] Progressive web app
- [ ] Mobile-specific features
- [ ] Cross-device testing

### Sprint 10: Launch and Feedback
- [ ] Production deployment
- [ ] Analytics implementation
- [ ] Admin dashboard
- [ ] User feedback system

## 🚦 Getting Started

### Prerequisites
- Go 1.24 or higher
- Node.js 18+ and npm/yarn
- PostgreSQL 13+ (when database is implemented)
- Google Cloud Console project (for OAuth credentials)
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/jointrip.git
   cd jointrip
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Install React dependencies**
   ```bash
   cd web
   npm install
   cd ..
   ```

4. **Set up Google OAuth credentials**
   ```bash
   # Create .env file with your Google OAuth credentials
   cp .env.example .env
   # Edit .env with your Google Client ID and Secret
   ```

5. **Build React frontend**
   ```bash
   cd web
   npm run build
   cd ..
   ```

6. **Run the application**
   ```bash
   go run main.go
   ```

### Development Setup

1. **Install development tools**
   ```bash
   # Install air for hot reloading (optional)
   go install github.com/cosmtrek/air@latest

   # Install golangci-lint for code quality
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

2. **Development mode (with hot reloading)**
   ```bash
   # Terminal 1: Start React dev server
   cd web
   npm run dev

   # Terminal 2: Start Go server with air
   air
   ```

3. **Production build and run**
   ```bash
   # Build React for production
   cd web
   npm run build
   cd ..

   # Run Go server with embedded React
   go run main.go
   ```

4. **Run tests**
   ```bash
   # Go tests
   go test ./...

   # React tests
   cd web
   npm test
   ```

## 📊 Database Design

The application uses a comprehensive entity-relationship model with the following core entities:

- **User**: User accounts and profiles
- **Trip**: Travel opportunities and details
- **TripParticipant**: User participation in trips
- **Message**: Direct communication between users
- **Expense**: Shared expense tracking
- **ExpenseShare**: Expense splitting details
- **UserRating**: Peer review system
- **Notification**: System notifications

For detailed database design, see [Entity Relationship Model](docs/entity_relationship_model.md).

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Follow Go best practices and conventions
- Write comprehensive tests for new features
- Update documentation for significant changes
- Use meaningful commit messages
- Ensure code passes linting and tests

## 📝 API Documentation

API documentation will be available once the core endpoints are implemented. The API follows RESTful principles with the following base structure:

```
GET    /api/v1/users          # List users
POST   /api/v1/users          # Create user
GET    /api/v1/users/:id      # Get user details
PUT    /api/v1/users/:id      # Update user
DELETE /api/v1/users/:id      # Delete user

GET    /api/v1/trips          # List trips
POST   /api/v1/trips          # Create trip
GET    /api/v1/trips/:id      # Get trip details
PUT    /api/v1/trips/:id      # Update trip
DELETE /api/v1/trips/:id      # Delete trip

# Additional endpoints for messages, expenses, ratings, etc.
```

## 🔒 Security

Security is a top priority for JoinTrip. We implement:

- **Authentication**: Google OAuth 2.0 with JWT session management
- **Authorization**: Role-based access control for different user types
- **Data Protection**: Input validation and SQL injection prevention (sqlx)
- **Privacy**: Granular privacy controls and data anonymization options
- **Audit Trail**: Comprehensive logging for security monitoring
- **OAuth Security**: Secure token handling and refresh mechanisms

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Team

- **Development Team**: [Your Team Information]
- **Project Manager**: [PM Information]
- **Designer**: [Designer Information]

## 📞 Support

For support and questions:
- Create an issue in this repository
- Contact the development team
- Check the documentation in the `docs/` folder

## 🙏 Acknowledgments

- Thanks to all contributors and beta testers
- Inspired by the need for better travel companion platforms
- Built with love for the travel community

---

**Happy Traveling! 🌟**

*JoinTrip - Connecting travelers, sharing experiences, creating memories.*
