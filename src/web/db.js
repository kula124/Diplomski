import { createClient } from '@supabase/supabase-js'

const supabaseUrl = 'https://amsvkdzwcxmodmgaktiq.supabase.co'
const supabaseKey = process.env.SUPABASE_KEY

const supabase = createClient(supabaseUrl, supabaseKey)

const auth = async (email, password) => {
  console.log("AYY")
  let { user, session, error } = await supabase.auth.signIn({
    email,
    password
  })
  if (error) {
    console.error(error)
    throw error
  }
  return { user, session }
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

const userInfo = async token => supabase.auth.api.getUser(token)

module.exports = {
  auth,
  signUp,
  logout,
  userInfo
}
