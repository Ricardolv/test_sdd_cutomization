export const messages = {
  'auth.name_required': 'Name is required',
  'auth.email_required': 'Email is required',
  'auth.email_duplicate': 'This email is already registered',
  'auth.password_required': 'Password is required',
  'auth.password_too_short': 'Password must be at least 6 characters',
  'auth.user_not_found': 'User not found',
  'auth.register_success': 'Account created successfully!',
  'auth.login_coming_soon': 'Login coming soon',
  'join.title_register': 'Create account',
  'join.title_login': 'Sign in',
  'join.switch_to_login': 'Already have an account? Sign in',
  'join.switch_to_register': "Don't have an account? Create account",
  'join.name_label': 'Name',
  'join.name_placeholder': 'Your name',
  'join.email_label': 'Email',
  'join.email_placeholder': 'you@email.com',
  'join.password_label': 'Password',
  'join.password_placeholder': 'Your password',
  'join.register_button': 'Create account',
  'join.login_button': 'Sign in',
} as const

export type MessageKey = keyof typeof messages
