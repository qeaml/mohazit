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
# yes, i know that this doesn't look that good; i'm working on it
if {num} = {guess}
    # that was a good guess!
    say Congratulations! You guessed the number!
else
    # well, better luck next time
    say Uh-oh! That's not correct. The number was: {num}
end
```

example with labels:

```rb
label greater
    # local variable, will be deleted after the 'end'
    local what = greater
    say Woah!
    say That number is {what} than 10!
    goto bye
end

label lower
    local what = lower
    say Well,
    say That number is {what} than or equal to 10.
    goto bye
end

# global variables, will never be deleted
set num = ask-number
set target = value 10
if {num} > {target}
    goto greater
else
    goto lower
end
# you use the same terminator for both `if`, `else` and `label` :)

label
    say Goodbye!
    exit
end
```

you aren't stuck with only `if` though:

```rb
# bruh
if this
else
    say woop
end
# ==== INTRODUCING: THE UNLESS STATEMENT ====
unless this = that
    say woop
end
# you can use else here too
unless this = that
    say woop
else
    say how??????
end
```
