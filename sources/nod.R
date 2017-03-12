
nodding <- function(angle, nod) {
  if (nod >= -1 && nod <= 1) {
    return (angle * nod)
  }
  sign <- 1
  if (nod < 0) {
    sign <- -1
    nod <- -nod
  }
  if (angle < 0) {
    return (sign * -0.5 * pi * (1.0 - (1.0 + angle*2.0/pi)^nod))
  }
  return   (sign *  0.5 * pi * (1.0 - (1.0 - angle*2.0/pi)^nod))
}

x <- ((0:100 / 50) - 1) * pi / 2
y <- x

nods <- c(4, 3, 2, 1, 0.5, -.9, -4)
colors <- c("black", "blue", "green", "red", "purple", "brown", "grey")


plot(x, y, type="n")

cl <- 1
for (nod in nods) {
  p <- 1
  for (xx in x) {
    y[p] <- nodding(xx, nod)
    p <- p+1
  }
  lines(x, y, type="l", col=colors[cl], lwd=2)
  cl <- cl+1
}

legend(-pi/2, .7, legend=nods, lwd=2, col=colors)
