import httpClient from '../utils/httpClient'

export default {
  // 任务列表
  list(query, callback) {
    httpClient.batchGet([
      {
        uri: '/task',
        params: query
      },
      {
        uri: '/executor/all'
      }
    ], callback)
  },

  detail(id, callback) {
    httpClient.batchGet([
      {
        uri: `/task/${id}`
      },
      {
        uri: '/executor/all'
      }
    ], callback)
  },

  update(data, callback) {
    httpClient.post('/task/save', data, callback)
  },

  remove(id, callback) {
    httpClient.post(`/task/delete/${id}`, {}, callback)
  },

  enable(id, callback) {
    httpClient.post(`/task/enable/${id}`, {}, callback)
  },

  disable(id, callback) {
    httpClient.post(`/task/disable/${id}`, {}, callback)
  },

  run(id, callback) {
    httpClient.get(`/task/run/${id}`, {}, callback)
  }
}
