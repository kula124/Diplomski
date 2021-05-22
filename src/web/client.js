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

module.exports = {
  auth,
  signUp,
  logout,
  userInfo,
  searchByHash
}
