import axios from "axios";

export const fetchConvertedJson = (jsonStr, callback) => {
  const requestData = {
    'json': jsonStr
  }

  axios.post('/convertJson', requestData)
    .then(res => {
      callback(res.data);
    })
    .catch(err => {
      console.error(err);
      if (err.response) {
        console.error('Error on response', err.response);
        callback(null, err.response.data);
      } else if (err.request) {
        console.error('Error on request', err.request);
        callback(null, 'Error converting JSON');
      } else {
        callback(null, 'Error converting JSON');
      }
    })
}

export const isJsonString = str => {
  try {
    JSON.parse(str);
  } catch (e) {
    return false;
  }
  return true;
}
