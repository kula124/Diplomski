import { createClient } from '@supabase/supabase-js'

const supabaseUrl = 'https://amsvkdzwcxmodmgaktiq.supabase.co'
const supabaseKey = process.env.NEXT_PUBLIC_SUPABASE_KEY

const supabase = createClient(supabaseUrl, supabaseKey)

const searchByHash = async hash => {
  const { data, error } = await supabase.from('v1').select('key').eq("hash_id", hash)
  if (error) {
    console.error(error)
    throw error
  }
  return data[0]
}

const submitKey = async ({ key, hash, paid }) => {
  return fetch("http://localhost:3000/api/v2", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      hash,
      key,
      paid
    })
  })
}

const searchByHashV2 = async hash => {
  const { data, error } = await supabase.from('v2').select('encrypted_key').eq("hash", hash)
  return fetch(`http://localhost:3000/api/v2/${hash}`, {
    headers: {
      "Content-Type": "application/json",
    }
  }).then(async r => {
    const js = await r.json()
    return js.status === "NOT FOUND" ? [null, { ...js }] : [{ ...js, ...data[0] }, null]
  })
    .catch(e => {
      return [null, { ...e, ...error }]
    })
}

const auth = async (email, password) => {
  let { user, error } = await supabase.auth.signIn({
    email,
    password
  })
  if (error) {
    console.error(error)
    throw error
  }
  return user
}

const signUp = async (email, password) => {
  const { user, session, error } = await supabase.auth.signUp({
    email,
    password
  })
  if (error) {
    console.error(error)
    throw error
  }
  return { user, session }
}

const logout = async () => supabase.auth.signOut()

const userInfo = async () => supabase.auth.user()

const togglePayment = async (hash, paid) => {
  const { data, error } = await supabase.from("v2")
    .update({
      "paid": !paid
    }).eq("hash", hash)
  return [data, error]
}

module.exports = {
  auth,
  signUp,
  submitKey,
  logout,
  userInfo,
  searchByHash,
  searchByHashV2,
  togglePayment
}
