import httpClient from '../utils/httpClient'

export default {
  // 任务列表
  list(query, callback) {
    httpClient.batchGet([
      {
        uri: '/job',
        params: query
      },
      {
        uri: '/executor'
      }
    ], callback)
  },

  detail(id, callback) {
    httpClient.batchGet([
      {
        uri: `/task/${id}`
      },
      {
        uri: '/executor'
      }
    ], callback)
  },

  update(data, callback) {
    httpClient.post('/job/save', data, callback)
  },

  remove(id, callback) {
    httpClient.post(`/job/delete/${id}`, {}, callback)
  },

  enable(id, callback) {
    httpClient.post(`/job/enable/${id}`, {}, callback)
  },

  disable(id, callback) {
    httpClient.post(`/job/disable/${id}`, {}, callback)
  },

  run(id, callback) {
    httpClient.get(`/job/run/${id}`, {}, callback)
  }
}
