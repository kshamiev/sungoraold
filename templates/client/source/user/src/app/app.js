angular.module('zegotaUsr', [
  'zegota.core',
  'zegota.ui',
  'zegota.recovery',
  'zegota.registration',
  'zegota.profile',
  'templates.user'//,
  //'ngMockE2E'
])
.config(function($interpolateProvider){
    $interpolateProvider.startSymbol('{[{').endSymbol('}]}');
  }
)
//.run(function($httpBackend) {

  //$httpBackend.whenPOST('/api/v1.0/session').respond(function(method,url,data) {
    //var icomingData, userWithCaptcha, userWithoutCaptcha, fakeUser, responseData;

    //userWithCaptcha = {
      //Login: 'alexp',
      //Password: '123123',
      //Captcha: '583191'
    //};

    //userWithoutCaptcha = {
      //Login: 'alexp',
      //Password: '123123',
    //};

    //incomingData = angular.fromJson(data);

    //for (prop in incomingData) {
      //if (incomingData.hasOwnProperty('Captcha')) {
        //fakeUser = userWithCaptcha;
      //} else {
        //fakeUser = userWithoutCaptcha;
      //}
    //}

    //if (JSON.stringify(fakeUser) === JSON.stringify(incomingData)) {
      //responseData = {
        //ErrorCode: 0,
        //ErrorMessage: '',
        //Token: 'a468d9e24f81fe38919d31a34e629d844d899be0470824afdb656b2d3899e5b5'
      //};
      //return [200, responseData];
    //} else {
      //responseData = {
        //ErrorCode: 2,
        //ErrorMessage: 'Ошибка,  логин или пароль указан не верно',
        //Token: ''
      //};
      //return [404, responseData];
    //}
  //});

  //$httpBackend.whenGET('/api/v1.0/session/')
    //.respond({
      //ErrorCode: 0,
      //ErrorMessage: '',
      //Token: 'a468d9e24f81fe38919d31a34e629d844d899be0470824afdb656b2d3899e5b5'
    //});

  //$httpBackend.whenGET('/api/v1.0/session/a468d9e24f81fe38919d31a34e629d844d899be0470824afdb656b2d3899e5b5')
    //.respond({
      //ErrorCode: 0,
      //ErrorMessage: ''
  //});

  //$httpBackend.whenGET('/api/v1.0/currentuser/a468d9e24f81fe38919d31a34e629d844d899be0470824afdb656b2d3899e5b5')
    //.respond({
      //Login: "alexp",
      //LastName: "p",
      //Name: "alex"
  //});

  //$httpBackend.whenDELETE('/api/v1.0/session/a468d9e24f81fe38919d31a34e629d844d899be0470824afdb656b2d3899e5b5')
    //.respond(function () {
      //return [200, {}, {}];
  //});

  //$httpBackend.whenGET('/api/v1.0/captcha/native')
    //.respond(function () {
      //var responseData, captchaImage;
      //captchaImage = "data:image/jpeg;base64,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";
      //responseData = {
        //ErrorCode: 0,
        //ErrorMessage: '',
        //CaptchaImage: captchaImage
      //};
      //return [200, responseData, {}];
  //});
//})
;
