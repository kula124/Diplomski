import React, { useState } from 'react'

import api from '../../client'

const v1 = () => {
  const [search, setSearch] = useState("")
  const [key, setKey] = useState()

  const keyDown = async e => {
    if (e.key === "Enter") {
      const res = await api.searchByHash(search)
      console.log(res)
      if (res) {
        return setKey(res)
      }
      setKey(false)
    }
  }

  return (
    <main className="bg-gray-800 h-screen w-screen flex flex-col justify-center items-center" onKeyPress={keyDown}>
      <section className="bg-gray-800 flex flex-col h-20 text-center w-2/6 text-teal text-lg">
        <span>Search the hash</span>
        <input type="text" className="placeholder-white text-center bg-gray-800 outline-none w-1/1 border-b border-teal"
          onChange={e => setSearch(e.target.value)} value={search} />
      </section>
      {key &&
        <span className="bg-gray-800 outline-none w-max text-center text-green-500 text-xl">
          {key.key}
        </span>
      }
    </main>
  )
}

export default v1