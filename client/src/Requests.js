export function sessionState() {
  const myHeaders = new Headers();
  myHeaders.append("x-user-uuid", "b38e6e7e-47c3-4256-ab7e-e513ba13ccdf");

  const requestOptions = {
    method: 'GET',
    headers: myHeaders,
    redirect: 'follow'
  };

  return fetch("http://0.0.0.0:8013/v1/session-status", requestOptions)
    .then(response => response.json())
    .catch(error => {
      console.log('error', error);
      return 'error'
    });
}
