import React from 'react'
import { cn } from '@/lib/utils'

interface ButtonGroupProps {
  children: React.ReactNode
  className?: string
  spacing?: 'tight' | 'normal' | 'loose'
  direction?: 'horizontal' | 'vertical'
}

const ButtonGroup: React.FC<ButtonGroupProps> = ({
  children,
  className,
  spacing = 'normal',
  direction = 'horizontal'
}) => {
  const spacingClasses = {
    tight: direction === 'horizontal' ? 'space-x-1' : 'space-y-1',
    normal: direction === 'horizontal' ? 'space-x-2' : 'space-y-2',
    loose: direction === 'horizontal' ? 'space-x-3' : 'space-y-3'
  }

  const directionClasses = {
    horizontal: 'flex flex-row',
    vertical: 'flex flex-col'
  }

  return (
    <div className={cn(
      directionClasses[direction],
      spacingClasses[spacing],
      className
    )}>
      {children}
    </div>
  )
}

export { ButtonGroup }
