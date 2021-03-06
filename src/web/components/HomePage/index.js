import React, { useState, useEffect } from 'react'
import { useRouter } from 'next/router'
import db from '../../client'

const textCss = 'cursor-pointer hover:bg-gray-700 transition-all shadow-2x1t duration-300 px-11 py-4 text-5xl text-teal flex flex-col justify-center items-center rounded border border-teal'

const Home = () => {
  const router = useRouter()
  return (
    <main className="flex bg-gray-800 h-screen flex-col justify-center items-center">
      <ul className="flex flex-col justify-evenly h-4/6 text-left" >
        <li className={textCss} onClick={() => { router.push('/v1') }}>
          <span>V1: Symmetric</span>
        </li>
        <li className={textCss} onClick={() => { router.push('/v2') }}>
          V2 & V3: Public RSA/Symmetric
        </li>
        <li className={textCss} onClick={() => { router.push('/keySubmit') }}>
          Submit key/hash pair
        </li>
        <li onClick={() => { db.logout(); setTimeout(() => router.reload()), 500 }}
          className='cursor-pointer hover:bg-gray-700 transition-all shadow-2x1 duration-300 px-11 py-4 text-5xl text-red-500 flex flex-col justify-center items-center rounded border border-red-500'>
          Logout
        </li>
      </ul>
    </main>
  )
}

export default Home