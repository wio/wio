function(write_sep)
    message(STATUS "#===========================================================================#")
endfunction()

function(info MESSAGE)
    message(STATUS ${MESSAGE})
endfunction()

function(fatal MESSAGE)
    message(FATAL_ERROR ${MESSAGE})
endfunction()

function(warning MESSAGE)
    message(WARNING ${MESSAGE})
endfunction()
