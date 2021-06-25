#
# Boost Software License - Version 1.0 - August 17th, 2003
#
# Permission is hereby granted, free of charge, to any person or organization
# obtaining a copy of the software and accompanying documentation covered by
# this license (the "Software") to use, reproduce, display, distribute,
# execute, and transmit the Software, and to prepare derivative works of the
# Software, and to permit third-parties to whom the Software is furnished to
# do so, all subject to the following:
# 
# The copyright notices in the Software and this entire statement, including
# the above license grant, this restriction and the following disclaimer,
# must be included in all copies of the Software, in whole or in part, and
# all derivative works of the Software, unless such copies or derivative
# works are solely in the form of machine-executable object code generated by
# a source language processor.
# 
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE, TITLE AND NON-INFRINGEMENT. IN NO EVENT
# SHALL THE COPYRIGHT HOLDERS OR ANYONE DISTRIBUTING THE SOFTWARE BE LIABLE
# FOR ANY DAMAGES OR OTHER LIABILITY, WHETHER IN CONTRACT, TORT OR OTHERWISE,
# ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
# DEALINGS IN THE SOFTWARE.
#
# 2012-01-31, Lars Bilke
# - Enable Code Coverage
#
# 2013-09-17, Joakim Söderberg
# - Added support for Clang.
# - Some additional usage instructions.
#
# USAGE:
# 1. Copy this file into your cmake modules path.
#
# 2. Add the following line to your CMakeLists.txt:
#      INCLUDE(CodeCoverage)
#
# 3. Set compiler flags to turn off optimization and enable coverage: 
#    SET(CMAKE_CXX_FLAGS "-g -O0 -fprofile-arcs -ftest-coverage")
#	 SET(CMAKE_C_FLAGS "-g -O0 -fprofile-arcs -ftest-coverage")
#  
# 3. Use the function SETUP_TARGET_FOR_COVERAGE to create a custom make target
#    which runs your test executable and produces a lcov code coverage report:
#    Example:
#	 SETUP_TARGET_FOR_COVERAGE(
#				my_coverage_target  # Name for custom target.
#				test_driver         # Name of the test driver executable that runs the tests.
#									# NOTE! This should always have a ZERO as exit code
#									# otherwise the coverage generation will not complete.
#				coverage            # Name of output directory.
#				)
#
# 4. Build a Debug build:
#	 cmake -DCMAKE_BUILD_TYPE=Debug ..
#	 make
#	 make my_coverage_target
#
#

# Check prereqs
FIND_PROGRAM( GCOV_PATH gcov )
FIND_PROGRAM( LCOV_PATH lcov )
FIND_PROGRAM( GENHTML_PATH genhtml )
FIND_PROGRAM( GCOVR_PATH gcovr PATHS ${CMAKE_SOURCE_DIR}/contrib/obs-studio/tests)

IF(NOT GCOV_PATH)
	MESSAGE(FATAL_ERROR "gcov not found! Aborting...")
ENDIF() # NOT GCOV_PATH

IF(NOT (CMAKE_COMPILER_IS_GNUCXX OR CMAKE_COMPILER_IS_GNUCC))
	# Clang version 3.0.0 and greater now supports gcov as well.
	MESSAGE(WARNING "Compiler is not GNU gcc! Clang Version 3.0.0 and greater supports gcov as well, but older versions don't.")
	
	IF(NOT ("${CMAKE_CXX_COMPILER_ID}" STREQUAL "Clang" OR "${CMAKE_C_COMPILER_ID}" STREQUAL "Clang"))
		MESSAGE(FATAL_ERROR "Compiler is not GNU gcc or Clang! Aborting...")
	ENDIF()
ENDIF() # NOT CMAKE_COMPILER_IS_GNUCXX

IF ( NOT CMAKE_BUILD_TYPE STREQUAL "Debug" )
  MESSAGE( WARNING "Code coverage results with an optimized (non-Debug) build may be misleading" )
ENDIF() # NOT CMAKE_BUILD_TYPE STREQUAL "Debug"


# Param _targetname     The name of new the custom make target
# Param _outputname     lcov output is generated as _outputname.info
#                       HTML report is generated in _outputname/index.html
# Param _testrunner     The name of the target which runs the tests.
#						MUST return ZERO always, even on errors. 
#						If not, no coverage report will be created!
# Optional fourth parameter is passed as arguments to _testrunner
#   Pass them in list form, e.g.: "-j;2" for -j 2
FUNCTION(SETUP_TARGET_FOR_COVERAGE _targetname _outputname _testrunner)

	IF(NOT LCOV_PATH)
		MESSAGE(FATAL_ERROR "lcov not found! Aborting...")
	ENDIF() # NOT LCOV_PATH

	IF(NOT GENHTML_PATH)
		MESSAGE(FATAL_ERROR "genhtml not found! Aborting...")
	ENDIF() # NOT GENHTML_PATH

	# Setup target
	ADD_CUSTOM_TARGET(${_targetname}
		
		# Cleanup lcov
		${LCOV_PATH} --directory . --zerocounters
		
		# Run tests
		COMMAND ${_testrunner} ${ARGV3}
		
		# Capturing lcov counters and generating report
		COMMAND ${LCOV_PATH} --directory . --capture --output-file ${_outputname}.info
		COMMAND ${LCOV_PATH} --remove ${_outputname}.info 'tests/*' '/usr/*' --output-file ${_outputname}.info.cleaned
		COMMAND ${GENHTML_PATH} -o ${_outputname} ${_outputname}.info.cleaned
		COMMAND ${CMAKE_COMMAND} -E remove ${_outputname}.info ${_outputname}.info.cleaned
		
		WORKING_DIRECTORY ${CMAKE_BINARY_DIR}
		COMMENT "Resetting code coverage counters to zero.\nProcessing code coverage counters and generating report."
	)
	
	# Show info where to find the report
	ADD_CUSTOM_COMMAND(TARGET ${_targetname} POST_BUILD
		COMMAND ;
		COMMENT "Open ./${_outputname}/index.html in your browser to view the coverage report."
	)

ENDFUNCTION() # SETUP_TARGET_FOR_COVERAGE

# Param _targetname     The name of new the custom make target
# Param _testrunner     The name of the target which runs the tests
# Param _outputname     cobertura output is generated as _outputname.xml
# Optional fourth parameter is passed as arguments to _testrunner
#   Pass them in list form, e.g.: "-j;2" for -j 2
FUNCTION(SETUP_TARGET_FOR_COVERAGE_COBERTURA _targetname _testrunner _outputname)

	IF(NOT PYTHON_EXECUTABLE)
		MESSAGE(FATAL_ERROR "Python not found! Aborting...")
	ENDIF() # NOT PYTHON_EXECUTABLE

	IF(NOT GCOVR_PATH)
		MESSAGE(FATAL_ERROR "gcovr not found! Aborting...")
	ENDIF() # NOT GCOVR_PATH

	ADD_CUSTOM_TARGET(${_targetname}

		# Run tests
		${_testrunner} ${ARGV3}

		# Running gcovr
		COMMAND ${GCOVR_PATH} -x -r ${CMAKE_SOURCE_DIR}/contrib/obs-studio -e '${CMAKE_SOURCE_DIR}/contrib/obs-studio/tests/'  -o ${_outputname}.xml
		WORKING_DIRECTORY ${CMAKE_BINARY_DIR}
		COMMENT "Running gcovr to produce Cobertura code coverage report."
	)

	# Show info where to find the report
	ADD_CUSTOM_COMMAND(TARGET ${_targetname} POST_BUILD
		COMMAND ;
		COMMENT "Cobertura code coverage report saved in ${_outputname}.xml."
	)

ENDFUNCTION() # SETUP_TARGET_FOR_COVERAGE_COBERTURA

