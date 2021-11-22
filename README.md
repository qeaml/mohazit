# mohazit

dead simple scripting

the entirety of the mohazit scripting mini-language is implemented in pure Go

here's an example of what mohazit code (will) look like:

```rb
# number up to (including) ten
set num = limited-rng 10
# ask the user for a number
set guess = ask-number
# compare
if `num` (equals) `guess`
    # that was a good guess!
    say Congratulations! You guessed the number!
else
    # well, better luck next time
    say Uh-oh! That's not correct. The number was:
    # could add templating or whatever it's called, but this is realistically
    # what it will look like in the first release
    say-var num
end
```

example with labels:

```rb
label greater
    say Woah!
    say That number is greater than 10!
    goto bye
end

label lower
    say Well,
    say That number is lower than or equal to 10.
    goto bye
end

set num ask-number
if `num` (greater) 10
    goto greater
else
    goto lower
end
# you use the same terminator for both `if`, `else` and `label` :)
```

you aren't stuck with only `if` though:

```rb
if this (equals) that
else
    say woop
end
# ==== INTRODUCING: THE UNLESS STATEMENT ====
unless this (equals) that
    say woop
end
# the code below won't run in the current version of mohazit, but it will eventually
unless this (equals) that
    say woop
else
    say how??????
end
```
