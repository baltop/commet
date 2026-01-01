# Commet Dashboard Test Plan

## Application Overview

Test plan for the Commet Dashboard - a Go Gin web application with HTMX + Alpine.js frontend. Covers authentication flows (login, register, logout), dashboard functionality (stats cards, charts), dark mode toggle, and responsive design.

## Test Scenarios

### 1. Authentication - Registration

**Seed:** `seed.spec.ts`

#### 1.1. should display registration form with all required fields

**File:** `tests/auth/register-form-display.spec.ts`

**Steps:**
  1. Navigate to /auth/register
  2. Verify name input field is visible
  3. Verify email input field is visible
  4. Verify password input field is visible
  5. Verify confirm password input field is visible
  6. Verify terms checkbox is visible
  7. Verify submit button is visible

**Expected Results:**
  - Registration form is displayed with all fields
  - Submit button is initially disabled

#### 1.2. should show validation error when submitting empty form

**File:** `tests/auth/register-empty-form.spec.ts`

**Steps:**
  1. Navigate to /auth/register
  2. Check the terms checkbox
  3. Click submit button
  4. Wait for error message

**Expected Results:**
  - Error message '모든 필드를 입력해주세요.' is displayed

#### 1.3. should show password strength indicator

**File:** `tests/auth/register-password-strength.spec.ts`

**Steps:**
  1. Navigate to /auth/register
  2. Type weak password 'abc'
  3. Verify weak indicator shown
  4. Type medium password 'abc123'
  5. Verify medium indicator shown
  6. Type strong password 'Abc123!@'
  7. Verify strong indicator shown

**Expected Results:**
  - Password strength indicator updates in real-time
  - Strength levels change from weak to strong based on password complexity

#### 1.4. should show password match indicator

**File:** `tests/auth/register-password-match.spec.ts`

**Steps:**
  1. Navigate to /auth/register
  2. Type password 'Test1234!'
  3. Type mismatched confirm password 'Different'
  4. Verify mismatch indicator (red X) is shown
  5. Clear and type matching confirm password 'Test1234!'
  6. Verify match indicator (green checkmark) is shown

**Expected Results:**
  - Red X indicator appears when passwords do not match
  - Green checkmark appears when passwords match

#### 1.5. should successfully register new user

**File:** `tests/auth/register-success.spec.ts`

**Steps:**
  1. Navigate to /auth/register
  2. Fill in name 'Test User'
  3. Fill in email with unique test email
  4. Fill in password 'TestPass123!'
  5. Fill in matching confirm password
  6. Check terms checkbox
  7. Click submit button
  8. Wait for redirect

**Expected Results:**
  - User is redirected to login page
  - Success message may be displayed

### 2. Authentication - Login

**Seed:** `seed.spec.ts`

#### 2.1. should display login form with all fields

**File:** `tests/auth/login-form-display.spec.ts`

**Steps:**
  1. Navigate to /auth/login
  2. Verify email input field is visible
  3. Verify password input field is visible
  4. Verify submit button is visible
  5. Verify link to registration page exists

**Expected Results:**
  - Login form is displayed with email and password fields
  - Submit button is visible and enabled

#### 2.2. should show error when submitting empty form

**File:** `tests/auth/login-empty-form.spec.ts`

**Steps:**
  1. Navigate to /auth/login
  2. Click submit button without entering credentials
  3. Wait for error message

**Expected Results:**
  - Error message '이메일과 비밀번호를 입력해주세요.' is displayed

#### 2.3. should show error for invalid credentials

**File:** `tests/auth/login-invalid-credentials.spec.ts`

**Steps:**
  1. Navigate to /auth/login
  2. Enter email 'invalid@test.com'
  3. Enter password 'wrongpassword'
  4. Click submit button
  5. Wait for error message

**Expected Results:**
  - Error message '이메일 또는 비밀번호가 올바르지 않습니다.' is displayed
  - Email field value is preserved

#### 2.4. should successfully login with valid credentials

**File:** `tests/auth/login-success.spec.ts`

**Steps:**
  1. Navigate to /auth/login
  2. Enter valid test user email
  3. Enter valid test user password
  4. Click submit button
  5. Wait for redirect to dashboard

**Expected Results:**
  - User is redirected to /dashboard
  - Dashboard page displays user's name in welcome message

### 3. Authentication - Logout

**Seed:** `seed.spec.ts`

#### 3.1. should logout from sidebar

**File:** `tests/auth/logout-sidebar.spec.ts`

**Steps:**
  1. Login with valid credentials
  2. Navigate to dashboard
  3. Click logout button in sidebar
  4. Wait for redirect

**Expected Results:**
  - User is redirected to login page
  - Attempting to access dashboard redirects back to login

#### 3.2. should redirect unauthenticated user to login

**File:** `tests/auth/protected-route.spec.ts`

**Steps:**
  1. Clear any existing auth cookies
  2. Navigate directly to /dashboard

**Expected Results:**
  - User is redirected to /auth/login

### 4. Dashboard - Statistics Cards

**Seed:** `seed.spec.ts`

#### 4.1. should display all four statistics cards

**File:** `tests/dashboard/stats-cards-display.spec.ts`

**Steps:**
  1. Login with valid credentials
  2. Navigate to dashboard
  3. Verify Total Users card is visible
  4. Verify Total Revenue card is visible
  5. Verify Total Orders card is visible
  6. Verify Conversion Rate card is visible

**Expected Results:**
  - All four statistics cards are displayed
  - Each card shows a value and growth percentage

#### 4.2. should display welcome message with user name

**File:** `tests/dashboard/welcome-message.spec.ts`

**Steps:**
  1. Login with valid credentials
  2. Navigate to dashboard
  3. Locate welcome message in header

**Expected Results:**
  - Welcome message displays '안녕하세요, [User Name]님!'

### 5. Dashboard - Charts

**Seed:** `seed.spec.ts`

#### 5.1. should load line chart for monthly sales

**File:** `tests/dashboard/chart-line.spec.ts`

**Steps:**
  1. Login with valid credentials
  2. Navigate to dashboard
  3. Wait for line chart container to load via HTMX
  4. Verify chart canvas is rendered

**Expected Results:**
  - Line chart titled '월별 매출 추이' is displayed
  - Chart canvas element is visible

#### 5.2. should load bar chart for category sales

**File:** `tests/dashboard/chart-bar.spec.ts`

**Steps:**
  1. Login with valid credentials
  2. Navigate to dashboard
  3. Wait for bar chart container to load via HTMX
  4. Verify chart canvas is rendered

**Expected Results:**
  - Bar chart titled '카테고리별 판매' is displayed
  - Chart canvas element is visible

#### 5.3. should load pie chart for traffic sources

**File:** `tests/dashboard/chart-pie.spec.ts`

**Steps:**
  1. Login with valid credentials
  2. Navigate to dashboard
  3. Wait for pie chart container to load via HTMX
  4. Verify chart canvas is rendered

**Expected Results:**
  - Pie chart titled '트래픽 소스' is displayed
  - Chart canvas element is visible

### 6. Dark Mode

**Seed:** `seed.spec.ts`

#### 6.1. should toggle dark mode on login page

**File:** `tests/theme/dark-mode-login.spec.ts`

**Steps:**
  1. Navigate to /auth/login
  2. Verify page is in light mode
  3. Click dark mode toggle button
  4. Verify page switches to dark mode

**Expected Results:**
  - Page background changes to dark colors
  - Text colors adjust for dark mode

#### 6.2. should persist dark mode preference

**File:** `tests/theme/dark-mode-persist.spec.ts`

**Steps:**
  1. Navigate to /auth/login
  2. Click dark mode toggle to enable dark mode
  3. Refresh the page
  4. Verify dark mode is still enabled

**Expected Results:**
  - Dark mode preference is saved to localStorage
  - Dark mode persists after page refresh

#### 6.3. should toggle dark mode on dashboard

**File:** `tests/theme/dark-mode-dashboard.spec.ts`

**Steps:**
  1. Login with valid credentials
  2. Navigate to dashboard
  3. Click dark mode toggle button
  4. Verify dashboard switches to dark mode
  5. Verify charts adapt to dark mode colors

**Expected Results:**
  - Dashboard theme changes to dark mode
  - All UI elements adapt to dark theme

### 7. Responsive Design

**Seed:** `seed.spec.ts`

#### 7.1. should show mobile menu on small screens

**File:** `tests/responsive/mobile-menu.spec.ts`

**Steps:**
  1. Login with valid credentials
  2. Set viewport to mobile size (375x667)
  3. Navigate to dashboard
  4. Verify hamburger menu button is visible
  5. Verify sidebar is hidden by default
  6. Click hamburger menu button
  7. Verify mobile sidebar slides in

**Expected Results:**
  - Hamburger menu is visible on mobile
  - Sidebar opens when hamburger is clicked
  - Dark overlay appears behind sidebar

#### 7.2. should close mobile sidebar when clicking outside

**File:** `tests/responsive/mobile-menu-close.spec.ts`

**Steps:**
  1. Login with valid credentials
  2. Set viewport to mobile size
  3. Navigate to dashboard
  4. Open mobile sidebar
  5. Click on the overlay outside sidebar

**Expected Results:**
  - Mobile sidebar closes
  - Overlay disappears
