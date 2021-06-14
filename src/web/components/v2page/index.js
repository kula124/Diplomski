import React, { useEffect, useState } from 'react'

import api from '../../client'

const v1 = () => {
  const [search, setSearch] = useState("")
  const [loading, setLoading] = useState(false)
  const [keyProps, setKeyProps] = useState()
  const [key, setKey] = useState()

  // const
  const approveClass = "outline-none bg-transparent w-max mt-10 mr-auto ml-auto resize-none hover:bg-green-400 text-green-600 font-semibold hover:text-white py-2 px-4 border border-green-400 hover:border-transparent rounded"
  const disproveClass = "outline-none bg-transparent w-max mt-10 mr-auto ml-auto resize-none hover:bg-red-500 text-red-700 font-semibold hover:text-white py-2 px-4 border border-red-500 hover:border-transparent rounded"
  const loadingClass = "outline-none bg-transparent w-max-30 mt-10 mr-auto ml-auto resize-none hover:bg-gray-500 hover:cursor-default text-gray-700 font-semibold hover:text-white py-2 px-12 border border-gray-500 hover:border-transparent rounded"

  const togglePayment = async hash => {
    setLoading(true)
    await api.togglePayment(hash, keyProps.status)
    await getKeyInfo(hash)
    setLoading(false)
  }

  const getKeyInfo = async hash => {
    const [keyProps, err] = await api.searchByHashV2(hash)
    if (err) {
      setKey(null)
      setKeyProps(null)
      return
    }
    const flag = keyProps.status === "PAID"
    // setKey(flag ? keyProps.key : keyProps.encrypted_key)
    if (keyProps) {
      return setKeyProps({ ...keyProps, status: flag })
    }
    setKeyProps(false)
  }

  useEffect(() => {
    if (!keyProps) {
      return
    }
    if (keyProps.status) {
      setKey(keyProps.key)
    } else {
      setKey(keyProps.encrypted_key)
    }
  }, [keyProps])

  const keyDown = e => {
    if (e.key === "Enter") {
      getKeyInfo(search)
    }
  }

  return (
    <main className="bg-gray-800 h-screen w-screen flex flex-col justify-center items-center" onKeyPress={keyDown}>
      <section className="bg-gray-800 flex flex-col h-20 text-center w-2/6 text-teal text-lg">
        <span>Search the hash</span>
        <input type="text" className="placeholder-white text-center bg-gray-800 outline-none w-1/1 border-b border-teal"
          onChange={e => setSearch(e.target.value)} value={search} />
      </section>
      {keyProps && <section className="flex flex-col justify-evenly text-teal">
        <label>{!keyProps.status ? "Encrypted key:" : "Key:"}</label>
        <textarea rows={keyProps.status ? "2" : "3"} style={{ width: "45ch" }} readOnly defaultValue={key}
          className="bg-gray-800 text-center max-w-lg scrollbar-thin text-green-400 resize-none outline-none scrollbar-thumb-teal scrollbar-track-transparent"
        />

        <label className={keyProps.status ? "text-green-500" : "text-red-500"}>
          Paid: {keyProps.status.toString()}
        </label>
      </section>
      }

      {key && <button style={{ cursor: loading ? "default" : "pointer" }} disabled={loading} className={loading ? loadingClass : (keyProps.status ? disproveClass : approveClass)} onClick={() => togglePayment(search)}>
        {loading ? "Loading..." : (keyProps.status ? "Disprove payment" : "Approve payment")}
      </button>}
    </main>
  )
}

export default v1