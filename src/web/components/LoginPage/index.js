import React, { useState, useEffect } from 'react'
import api from '../../client'
import { useRouter } from 'next/router'

const LoginPage = () => {
  const [user, setUser] = useState("")
  const [password, setPassword] = useState("")
  const router = useRouter()

  const keyDown = e => {
    if (e.key === "Enter")
      Submit(user, password)
  }

  const Submit = async (user, password) => {
    const { user: u, session } = await api.auth(user, password)
      .then(() => router.replace('/home'))
      .catch(err => {
        return console.error(err)
      })
    console.log(u, session)
  }


  return (
    <main className="h-screen bg-gray-900 w-screen flex items-center justify-center" onKeyDown={keyDown}>
      <div className="flex flex-col bg-gray-800 justify-evenly shadow-2x1t items-center w-1/6 h-2/6">
        <div className="flex flex-col justify-around h-1/8 text-center text-teal">
          <span>Username</span>
          <input className="bg-gray-800 text-indigo-400 focus:outline-none border-b-2 border-teal" type="text" onChange={e => setUser(e.target.value)} value={user} />
        </div>
        <div className="flex flex-col justify-around h-1/8 text-center text-teal">
          <span>Password</span>
          <input className="bg-gray-800 text-indigo-400 border-b-2 focus:outline-none border-teal" type="password" onChange={e => setPassword(e.target.value)} value={password} />
        </div>
      </div>
    </main>
  )
}

export default LoginPage