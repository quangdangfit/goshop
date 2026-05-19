type SSOProvider = 'google' | 'facebook' | 'apple'

/**
 * Brand glyphs rendered inline so we don't add another icon dependency.
 * Sized to drop into the SSO buttons below.
 */
function GoogleIcon() {
  return (
    <svg className="h-5 w-5" viewBox="0 0 48 48" aria-hidden="true">
      <path fill="#FFC107" d="M43.6 20.5H42V20H24v8h11.3c-1.6 4.7-6 8-11.3 8a12 12 0 1 1 7.9-21l5.7-5.7A20 20 0 1 0 24 44c11 0 20-9 20-20 0-1.2-.1-2.4-.4-3.5z"/>
      <path fill="#FF3D00" d="M6.3 14.7l6.6 4.8A12 12 0 0 1 24 12c3.1 0 5.9 1.2 8 3.1l5.7-5.7A20 20 0 0 0 6.3 14.7z"/>
      <path fill="#4CAF50" d="M24 44c5.2 0 10-2 13.6-5.2l-6.3-5.3A12 12 0 0 1 12.7 28l-6.6 5.1A20 20 0 0 0 24 44z"/>
      <path fill="#1976D2" d="M43.6 20.5H42V20H24v8h11.3a12 12 0 0 1-4.1 5.5l6.3 5.3c-.4.4 6.7-4.8 6.7-14.8 0-1.2-.1-2.4-.4-3.5z"/>
    </svg>
  )
}

function FacebookIcon() {
  return (
    <svg className="h-5 w-5" viewBox="0 0 24 24" aria-hidden="true">
      <path fill="#1877F2" d="M24 12a12 12 0 1 0-13.9 11.9v-8.4H7.1V12h3V9.4c0-3 1.8-4.7 4.5-4.7 1.3 0 2.7.2 2.7.2v3h-1.5c-1.5 0-2 .9-2 1.9V12h3.3l-.5 3.5h-2.8v8.4A12 12 0 0 0 24 12z"/>
    </svg>
  )
}

function AppleIcon() {
  return (
    <svg className="h-5 w-5" viewBox="0 0 24 24" aria-hidden="true">
      <path fill="currentColor" d="M16.4 12.7c0-2.4 2-3.6 2.1-3.6-1.1-1.7-2.9-1.9-3.5-1.9-1.5-.2-2.9.9-3.7.9-.8 0-1.9-.9-3.2-.9-1.6 0-3.2 1-4 2.4-1.7 3-.4 7.4 1.2 9.8.8 1.2 1.8 2.5 3.1 2.4 1.2 0 1.7-.8 3.2-.8s1.9.8 3.2.8c1.3 0 2.2-1.2 3-2.4.9-1.4 1.3-2.7 1.3-2.8 0 0-2.5-1-2.7-3.9zM14 5.4c.7-.8 1.1-1.9 1-3-1 0-2.1.6-2.8 1.4-.6.7-1.2 1.8-1 2.9 1.1.1 2.2-.6 2.8-1.3z"/>
    </svg>
  )
}

const PROVIDERS: { id: SSOProvider; label: string; Icon: () => JSX.Element }[] = [
  { id: 'google', label: 'Google', Icon: GoogleIcon },
  { id: 'facebook', label: 'Facebook', Icon: FacebookIcon },
  { id: 'apple', label: 'Apple', Icon: AppleIcon },
]

/**
 * SSOButtons renders the three upstream-IdP buttons plus a fallback "Sign in
 * with SSO" link that hits Authentik's default flow. Each button is a plain
 * anchor so the browser follows the 302 chain through Authentik back to
 * /auth/callback.
 */
export default function SSOButtons() {
  return (
    <div>
      <div className="grid grid-cols-3 gap-2">
        {PROVIDERS.map(({ id, label, Icon }) => (
          <a
            key={id}
            href={`/api/v1/auth/login?provider=${id}`}
            aria-label={`Continue with ${label}`}
            className="flex items-center justify-center gap-2 py-2 px-3 border border-gray-200 rounded-md text-sm text-gray-700 hover:bg-gray-50"
          >
            <Icon />
            <span className="hidden sm:inline">{label}</span>
          </a>
        ))}
      </div>

      <a
        href="/api/v1/auth/login"
        className="mt-3 block text-center text-sm text-primary-600 hover:text-primary-700"
      >
        Sign in with SSO
      </a>
    </div>
  )
}
