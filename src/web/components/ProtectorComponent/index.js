import { useRouter } from 'next/router'
import React, { useState, useEffect } from 'react'
import api from '../../client'

const Protector = ({ children, shouldBeLoggedIn, goTo, duration }) => {
  const router = useRouter()
  const [loading, setLoading] = useState(true)
  const [user, setUser] = useState(null)

  useEffect(() => {
    (async function () {
      setLoading(true)
      const u = await api.userInfo().catch(err => {
        console.error('Failed to get user', err)
        return false
      })
      // const u = false
      let isLoggedIn
      if (!u) {
        isLoggedIn = false
      } else {
        isLoggedIn = u.aud === "authenticated"
      }

      setUser(shouldBeLoggedIn === isLoggedIn)
      setLoading(false)
    })()
  }, [])

  if (loading) {
    return <main className="bg-gray-800 h-screen w-screen flex justify-center items-center">
      <p className="text-7xl text-teal">
        LOADING....
      </p>
    </main>
  }

  if (!user) {
    if (goTo) {
      setTimeout(() => router.replace(goTo), duration || 2000)
    }
    return <main className="bg-gray-800 h-screen w-screen flex justify-center items-center">
      <p className="text-7xl text-teal">
        {shouldBeLoggedIn ? "You are not logged in" : "You are logged in: redirecting..."}
      </p>
    </main>
  }

  return children
}

export default Protector