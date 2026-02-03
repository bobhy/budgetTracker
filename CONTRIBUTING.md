Welcome, stranger!

Like sudo, we trust you have been given the usual warnings about playing nice with others.

If you have any problems or suggestions, please open an issue.
If you have a fix, submit a PR!

## Tech Stack

Template for the app powered by https://github.com/rgehrsitz/wails2-svelte5-tailwind4-ts-vite.

This template combines the latest versions of powerful technologies:

- **[Wails v2.10.1](https://wails.io/)**: Build desktop applications using Go and web technologies
- **[Svelte v5.28.2](https://svelte.dev/)**: Cybernetically enhanced web apps with revolutionary reactivity
- **[Tailwind CSS v4.1.4](https://tailwindcss.com/)**: Utility-first CSS framework with new CSS-first configuration
- **[shadcn-svelte v1.0.6](https://shadcn-svelte.com/)**: Beautifully designed components built with Radix UI and Tailwind CSS
- **[TypeScript v5.8.3](https://www.typescriptlang.org/)**: JavaScript with syntax for types
- **[Vite v6.3.3](https://vitejs.dev/)**: Next generation frontend tooling for lightning-fast development

### Features

- **🎨 Complete UI Component Library**: 40+ pre-built, accessible components from shadcn-svelte
- **🌙 Dark Mode Support**: Built-in dark/light theme switching with proper color variables
- **⚡ Modern Development**: Svelte 5's runes system with Tailwind 4's CSS-first configuration
- **🔧 Type Safety**: Full TypeScript support throughout the project
- **🚀 Fast Development**: Hot module replacement powered by Vite with @tailwindcss/vite plugin
- **📱 Responsive Design**: Mobile-first approach with Tailwind's responsive utilities
- **♿ Accessibility**: Components built with accessibility best practices
- **🎯 Cross-Platform**: Build for Windows, macOS, and Linux with a single codebase
- **🔥 Go Backend**: Leverage Go's performance and ecosystem for your application logic

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.23 or later)
- [Node.js](https://nodejs.org/) (version 16 or later)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)

### Quick Start

1. Clone this template:
```bash
git clone https://github.com/your-username/wails2-svelte5-tailwind4-ts-vite.git
cd wails2-svelte5-tailwind4-ts-vite
```

2. Install dependencies:
```bash
# Frontend dependencies are installed automatically by Wails
wails dev
```

### Development

To run in live development mode:

```bash
wails dev
```

This will:
- Start a Go backend server
- Launch a Vite development server with hot reload
- Open your application in a native window
- Enable access via browser at http://localhost:34115

For frontend-only development:
```bash
cd frontend
npm run dev
```

### Building

To build a production-ready distributable package:

```bash
wails build
```

## Using shadcn-svelte Components

Import and use components in your Svelte files:

```svelte
<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Card, CardContent, CardHeader, CardTitle } from "$lib/components/ui/card";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
</script>

<Card class="w-96">
  <CardHeader>
    <CardTitle>Login</CardTitle>
  </CardHeader>
  <CardContent class="space-y-4">
    <div class="space-y-2">
      <Label for="email">Email</Label>
      <Input id="email" type="email" placeholder="Enter your email" />
    </div>
    <Button class="w-full">Sign In</Button>
  </CardContent>
</Card>
```

## Dark Mode

Dark mode is automatically configured. Toggle between themes:

```svelte
<script lang="ts">
  import { toggleMode } from "mode-watcher";
</script>

<Button on:click={toggleMode}>Toggle Theme</Button>
```

## Project Structure

```
├── frontend/                    # Svelte frontend application
│   ├── src/
│   │   ├── lib/
│   │   │   ├── components/
│   │   │   │   └── ui/         # shadcn-svelte components
│   │   │   ├── utils.ts        # Utility functions
│   │   │   └── hooks/          # Custom Svelte hooks
│   │   ├── App.svelte          # Main application component
│   │   ├── main.ts             # Application entry point
│   │   └── style.css           # Global styles with Tailwind
│   ├── components.json         # shadcn-svelte configuration
│   ├── package.json            # Frontend dependencies
│   ├── tsconfig.json           # TypeScript configuration
│   └── vite.config.ts          # Vite configuration
├── app.go                      # Go application context
├── main.go                     # Go application entry point
├── wails.json                  # Wails configuration
└── build/                      # Build assets and configuration
```

## Configuration

### Tailwind CSS 4

This template uses Tailwind CSS 4 with the new CSS-first configuration approach. All theme variables are defined in `frontend/src/style.css`:

- CSS custom properties for colors
- Built-in dark mode support
- tw-animate-css for animations
- @tailwindcss/vite for optimal performance

### shadcn-svelte Components

Components are configured in `frontend/components.json` and installed in `frontend/src/lib/components/ui/`. Each component is:

- Fully customizable and owns its own code
- Built with accessibility in mind
- Styled with Tailwind CSS
- TypeScript ready

### Path Aliases

The following path aliases are configured:

- `$lib` → `frontend/src/lib`
- `$lib/components` → `frontend/src/lib/components`
- `$lib/components/ui` → `frontend/src/lib/components/ui`
- `$lib/utils` → `frontend/src/lib/utils`
- `$lib/hooks` → `frontend/src/lib/hooks`

## Adding New Components

To add additional shadcn-svelte components:

```bash
cd frontend
npx shadcn-svelte@latest add [component-name]
```

For example:
```bash
npx shadcn-svelte@latest add calendar
npx shadcn-svelte@latest add date-picker
```

## Customization

This template provides a solid foundation that's easy to extend:

### Adding Custom Styles
- Modify `frontend/src/style.css` for global styles
- Customize color schemes by updating CSS custom properties
- Add custom Tailwind utilities using the `@layer` directive

### Extending Components
- All shadcn-svelte components are in your codebase and fully customizable
- Create new components in `frontend/src/lib/components/`
- Follow the established patterns for consistency

### Go Backend Integration
- Add your application logic in Go files
- Use Wails context for frontend-backend communication
- Leverage Go's standard library and ecosystem

### Environment Configuration
- Configure different environments in `wails.json`
- Set up environment variables for different build targets
- Customize build flags and assets per platform

## Development Tips

### Hot Reload
- Changes to Svelte components reload instantly
- Go code changes trigger automatic recompilation
- CSS changes apply immediately with Vite HMR

### Debugging
- Use browser dev tools for frontend debugging
- Access frontend in browser mode: `http://localhost:34115`
- Use Go debugging tools for backend investigation

### Performance
- Vite handles optimal bundling and tree shaking
- Tailwind CSS purges unused styles automatically
- shadcn-svelte components are lightweight and performant

## Testing

The budgetTracker project includes comprehensive test coverage for both backend and frontend code.

### Backend Tests (Go)

The backend uses Go's built-in testing framework. Tests are located in files ending with `_test.go`.

#### Run All Backend Tests

```bash
go test ./...
```

#### Run Tests for a Specific Package

```bash
# Test configuration package
go test ./config

# Test models package
go test ./models
```

#### Run Tests with Verbose Output

```bash
go test -v ./...
```

#### Run Tests with Coverage

```bash
go test -cover ./...
```

**Available Backend Tests:**
- `config/config_test.go` - Tests for configuration management, environment variables, and command-line flags
- `models/generics_test.go` - Tests for GORM generic CRUD operations and pagination

### Frontend Tests

The frontend has two types of tests: unit tests (Vitest) and end-to-end tests (Playwright). Frontend tests are located in the `frontend/` directory.

#### Unit Tests (Vitest)

Unit tests verify individual component logic and utilities. They run in a JSDOM environment.

**Run All Unit Tests:**

```bash
cd frontend
npx vitest run
```

**Run Tests in Watch Mode (for development):**

> **Note:** In the fish shell, watch mode may hang. Use the `run` command above or switch to bash.

```bash
cd frontend
npx vitest
```

#### End-to-End Tests (Playwright)

Playwright tests verify the DataTable component's behavior in a real browser environment.

**Run All Playwright Tests:**

```bash
cd /home/bobhy/worktrees/dataTable
npx playwright test --reporter=list
```

### Running All Tests

To verify the entire project before committing changes:

```bash
# 1. Run backend tests
go test ./...

# 2. Run frontend unit tests
cd frontend
npx vitest run

# 3. Run frontend E2E tests (optional, in separate workspace)
cd /home/bobhy/worktrees/dataTable
npx playwright test --reporter=list
```

### Test Best Practices

- **Before Submitting PRs:** Always run both backend and frontend tests to ensure no regressions
- **Backend Tests:** Use table-driven tests for comprehensive coverage
- **Frontend Unit Tests:** Keep tests isolated and use JSDOM for fast execution
- **Playwright Tests:** Use for critical user journeys and complex interactions
- **CI/CD:** Consider running tests automatically on push/PR

## Browser Compatibility

This template supports modern browsers with:
- ES2020+ features
- CSS custom properties
- CSS Grid and Flexbox
- Modern JavaScript APIs

For broader compatibility, configure Vite's build target in `frontend/vite.config.ts`.

## License

This template is available under the MIT License.
