import request from './config'

export function login(username, password) {
    return request({
        url: '/login',
        method: 'post',
        data: {
            username,
            password
        }
    })
}

export function logout() {
    console.log('Calling logout API')
    return request({
        url: '/logout',
        method: 'post',
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
    })
} 