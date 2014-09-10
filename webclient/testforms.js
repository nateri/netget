/*
 $Id: testforms.js,v 1.6 2008/07/04 10:22:56 agoubard Exp $

 Copyright 2003-2008 Online Breedband B.V.
 See the COPYRIGHT file for redistribution and use restrictions.
*/

function doHandleSubmit(form)
{
   var elemInUrl = document.getElementById('in_url');
   var elemUrlText = document.getElementById('url_text');
   var elemFunction = document.getElementById('in_function');
   
   elemInUrl.value = elemUrlText.value;
   elemFunction.value = "ExecRemoteFile";
   
   if (document.pressed == 'magnet')
   {
      elemFunction.value = "ExecProtocol";
   }
   else if(document.pressed == 'file.torrent')
   {
      elemInUrl.value = elemInUrl.value + ".torrent";
   }
   
   return true;
}

// Displays the request URL with the colors and execute the request
function doRequest(form) {
   var elems = form.elements;
   var iframe = document.getElementById('xmlOutputFrame');
   var querySpan = document.getElementById('query');
   var requestParams = [];
   var formattedRequestString = '';
   var value, name, requestString;
   var i;
   
   doHandleSubmit(form);
   
   iframe.src = "about:blank";
   for (i = 0; i != elems.length; i++) {
      if (!(name = elems[i].name) || name == '_environment' || name == '_autofill') {
         continue;
      }

      if (elems[i].type == 'text' || elems[i].type == 'hidden' || elems[i].type == 'textarea') {
         value = elems[i].value;
      } else if (elems[i].type == 'select-one') {
         value = elems[i].options[elems[i].selectedIndex].value;
      }

      if (value) {
         if (name == '_action' || name == '_method' || name == '_target') {
            name = name.substring(1);
         }
         if (window.encodeURIComponent) {
            value = encodeURIComponent(value);
         } else {
            value = escape(value);
         }
         requestParams[requestParams.length] = name + '=' + value;
         if (formattedRequestString) {
            formattedRequestString += '&amp;';
         }

         if (name == '_function') {
            formattedRequestString += '<span class="functionparam">';
         } else {
            formattedRequestString += '<span class="param">';
         }

         formattedRequestString += '<span class="name">' + name + '</span>';
         formattedRequestString += '=<span class="value">' + value + '</span>';
         formattedRequestString += '</span>';
      }
   }

   requestString = form.action + '?' + requestParams.join('&');
   formattedRequestString = form.action + '?' + formattedRequestString;

   iframe.src = requestString;
   querySpan.innerHTML = formattedRequestString;
   return false;
}

function getCookie(name) {
   var start = document.cookie.indexOf(name + "=");
   var len = start + name.length + 1;
   if ((!start) && (name != document.cookie.substring(0, name.length))) {
      return null;
   }
   if (start == -1) { return null; }
   var end = document.cookie.indexOf(";", len);
   if (end == -1) { end = document.cookie.length; }
   return unescape(document.cookie.substring(len, end));
}

function setCookie(name, value, expires, path, domain, secure) {
   var today = new Date();
   today.setTime(today.getTime());
   if (expires) {
      expires = expires * 1000 * 60 * 60 * 24;
   }
   var expires_date = new Date(today.getTime() + (expires));
   document.cookie = name + "=" + escape(value) +
      ((expires) ? ";expires=" + expires_date.toGMTString() : "") + //expires.toGMTString()
      ((path) ? ";path=" + path : "") +
      ((domain) ? ";domain=" + domain : "") +
      ((secure) ? ";secure" : "");
}


function deleteCookie(name, path, domain) {
   if (getCookie(name)) { document.cookie = name + "=" +
      ((path) ? ";path=" + path : "") +
      ((domain) ? ";domain=" + domain : "") +
      ";expires=Thu, 01-Jan-1970 00:00:01 GMT";
   }
}


function setEnvCookie(form) {
   var env;
   if (form._environment.options) {
      var selIndex = form._environment.selectedIndex;
      env = form._environment.options[selIndex].text;
      setCookie("xins.env", env, "", "", "", "");
   } else { //if (form._environment.type == 'text') {
      env = form._environment.value;
      setCookie("xins.env", env, "", "", "", "");
   }
}


function selectEnv() {
   var i;
   // make sure that only pages with form and environment set selected value from the env cookie
   if (document.forms[0] && document.forms[0]._environment && document.forms[0]._environment.options) {
      var options = document.forms[0]._environment.options;
      env = getCookie("xins.env");
      for (i = 0; i != options.length; i++) {
         var option = options[i];
         if (env == options[i].text) {
            options.selectedIndex = i;
         }
      }
   } else if (document.forms[0] && document.forms[0]._environment && document.forms[0]._environment.type == 'text') {
      env = getCookie("xins.env");
      if (env != null && env != 'null') {
         document.forms[0]._environment.value = env;
      }
   }
}
