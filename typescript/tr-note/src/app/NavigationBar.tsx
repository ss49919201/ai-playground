"use client";

import Link from "next/link";
import { useSession, signOut } from "../lib/auth";
import { useRouter } from "next/navigation";

const navStyles = {
  navbar: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
    padding: "1rem 2rem",
    backgroundColor: "var(--gray-alpha-100)",
    borderBottom: "1px solid var(--gray-alpha-200)",
  },
  logo: {
    fontWeight: "bold",
    fontSize: "1.2rem",
  },
  nav: {
    display: "flex",
    gap: "1.5rem",
    alignItems: "center",
  },
  navLink: {
    fontWeight: 500,
    transition: "opacity 0.2s",
  },
  navLinkHover: {
    opacity: 0.8,
  },
  button: {
    padding: "0.5rem 1rem",
    borderRadius: "6px",
    border: "none",
    backgroundColor: "var(--foreground)",
    color: "var(--background)",
    fontWeight: 500,
    cursor: "pointer",
    transition: "background-color 0.2s",
  },
  buttonHover: {
    backgroundColor: "var(--button-primary-hover)",
  },
  userInfo: {
    display: "flex",
    alignItems: "center",
    gap: "0.75rem",
  },
  userName: {
    fontWeight: 500,
  },
};

export default function NavigationBar() {
  const { session, loading } = useSession();
  const router = useRouter();

  const handleSignOut = async () => {
    await signOut();
    router.refresh();
  };

  return (
    <header style={navStyles.navbar}>
      <Link href="/" style={navStyles.logo}>
        TR Note
      </Link>

      <nav style={navStyles.nav}>
        {!loading && (
          <>
            {session ? (
              <div style={navStyles.userInfo}>
                <span style={navStyles.userName}>{session.user?.name || "ユーザー"}</span>
                <button 
                  onClick={handleSignOut}
                  style={navStyles.button}
                  onMouseOver={(e) => 
                    Object.assign(e.currentTarget.style, navStyles.buttonHover)
                  }
                  onMouseOut={(e) => 
                    e.currentTarget.style.backgroundColor = "var(--foreground)"
                  }
                >
                  ログアウト
                </button>
              </div>
            ) : (
              <>
                <Link 
                  href="/login" 
                  style={navStyles.navLink}
                  onMouseOver={(e) => 
                    Object.assign(e.currentTarget.style, navStyles.navLinkHover)
                  }
                  onMouseOut={(e) => 
                    e.currentTarget.style.opacity = "1"
                  }
                >
                  ログイン
                </Link>
                <Link 
                  href="/signup" 
                  style={{
                    ...navStyles.button,
                    textDecoration: "none",
                  }}
                  onMouseOver={(e) => 
                    Object.assign(e.currentTarget.style, navStyles.buttonHover)
                  }
                  onMouseOut={(e) => 
                    e.currentTarget.style.backgroundColor = "var(--foreground)"
                  }
                >
                  新規登録
                </Link>
              </>
            )}
          </>
        )}
      </nav>
    </header>
  );
}