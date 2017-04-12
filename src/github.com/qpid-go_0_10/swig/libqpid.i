%module libqpid
%{
#include "libqpid.h"
%}



%include "std_string.i"
%include "libqpid.h"

//https://github.com/jsolmon/go-swig-exceptions
%include "exception.i"

%exception {
  try {
    $action
  } catch (const std::exception& e) {
    SWIG_exception(SWIG_ValueError, e.what());
  }
}
