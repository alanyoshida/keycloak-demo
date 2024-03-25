import React from 'react'
import { useNavigate } from 'react-router-dom'
import GetPets from './getPets'

const Home = (props) => {

  const { loggedIn, username } = props
  const navigate = useNavigate()

  const onButtonClick = () => {
    if (loggedIn) {
      localStorage.removeItem('user')
      localStorage.removeItem('token')
      props.setLoggedIn(false)
      navigate('/login')
    } else {
      navigate('/login')
    }
  }

  return (
    <div className="mainContainer">
      <div className={'titleContainer'}>
        <div>Welcome!</div>
      </div>
      <div>This is the home page.</div>
      <h3>Pet List</h3>
      <div>
      { loggedIn ? <GetPets/> : "" }
      </div>
      <div className={'buttonContainer'}>
        <input
          className={'inputButton'}
          type="button"
          onClick={onButtonClick}
          value={loggedIn ? 'Log out' : 'Log in'}
        />
        {loggedIn ? <div>Your Username is {username}</div> : <div />}
      </div>
    </div>
  )
}

export default Home