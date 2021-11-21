# (Magical) Magnetometer Calibration

This library provides hassle-free implicit real-time calibration of magnetometer for [portable devices](https://tinygo.org).  

Ideally, no manual calibration shall be needed. The library eventually finds some good enough solution that can then be extracted and stored on device flash to warm-up the same library on the next boot.  
_Flash manipulations are highly device specific, not covered here and shall be dealt separately._

The library can adapt to changing environment and provide stable low-error results on the go, something that is important for portable.

Internally, the library applies gradient descent to task of finding magnetometer calibration parameters.  
Optimisation for embedded use brings `float32` instead of standard `float64`.  
Calibration happens on background, continuously.


## Magnetometer calibration is a tedious task

Magnetometer sensor returns a vector along Earth magnetic field.  
Readings collected while rotating the sensor around shall produce points that lay on a sphere with its center in sensor.

Sensor though is affected by hard- and soft-iron effects.
These effects distort perfect sphere to something resembling an elongated melon, offset from origin too.

There is a [mathematical model that describes this distortion](https://www.vectornav.com/resources/inertial-navigation-primer/specifications--and--error-budgets/specs-hsicalibration) and a formula to bring raw readings back on perfect centered sphere.
The task of calibration is to find parameters for that formula, since every sensor is different and no two environments are the same.

First, we need to collect some raw data to work with.  
Next, one can try and solve the task analytically, in theory.  
In practice, due to noise in incoming data, it's more convenient to use classic gadient descent approach and find a solution that's _good enough_.

## Continuous magnetometer calibration on background

Library keeps current best state and applies calibration transformation to all incoming vectors.  
It also keeps a buffer of previous seen vectors for reference and use them to run gradient descent on.  
Vectors in the buffer are organised in quadrants and library tries to maintain diversity here, this helps with finding the best solution.

After processing an incoming vector, library may trigger re-calibration on background.  
This heppens when length of resulting vector differs from target value too much.  
If that happens and no calibration is already running on background, this new vector is pushed to the buffer (releasing some another vector) and a goroutine starts that tries to reconcile vectors in the buffer with calibration state and find a solution that gives minimal error.

This library is targeted for [embedded systems](https://tinygo.org) with low resources there number or cores is usually limited to just one and no real parallelism possible (cooperative runtime). Thus, to let other processes to run smoothly, the background task can be throttled.
