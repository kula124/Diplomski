import Head from 'next/head'

import LoginPage from '../components/LoginPage'
import Protector from '../components/ProtectorComponent'

export default function Home() {
  return (
    <div>
      <Head>
        <title>RW CnC</title>
        <meta name="description" content="cool page brah" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Protector shouldBeLoggedIn={false} goTo="/home">
        <LoginPage />
      </Protector>
    </div>
  )
}
