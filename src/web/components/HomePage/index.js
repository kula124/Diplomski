import React, { useState, useEffect } from 'react'
import { useRouter } from 'next/router'
const textCss = 'cursor-pointer hover:bg-gray-700 transition-all shadow-2x1t duration-300 px-11 py-4 text-5xl text-teal flex flex-col justify-center items-center rounded border border-teal'

const Home = () => {
  const router = useRouter()
  return (
    <main className="flex bg-gray-800 h-screen flex-col justify-center items-center">
      <ul className="flex flex-col justify-evenly h-4/6 text-left" >
        <li className={textCss} onClick={() => { router.push('/v1') }}>
          <span>V1: Symmetric</span>
        </li>
        <li className={textCss}>
          V2: Public/Symmetric
      </li>
        <li className={textCss}>
          V3: True Hybrid
      </li>
      </ul>
    </main>
  )
}

export default Home